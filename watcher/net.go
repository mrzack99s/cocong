package watcher

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/uuid"
	"github.com/mrzack99s/cocong/model"
	"github.com/mrzack99s/cocong/session"
	"github.com/mrzack99s/cocong/utils"
	"github.com/mrzack99s/cocong/vars"
)

func NetWatcher(ctx context.Context) {
	go func(ctx context.Context) {
		time.Sleep(5 * time.Second)
		for {
			select {
			case <-time.After(500 * time.Millisecond):
			case <-ctx.Done():
				return
			default:
				intIpv4, err := utils.GetSecureInterfaceIpv4Addr()
				if err != nil {
					panic(err)
				}

				handle, err := pcap.OpenLive(vars.Config.SecureInterface, 65535, true, pcap.BlockForever)
				if err != nil {
					panic(err)
				}
				if err := handle.SetBPFFilter(fmt.Sprintf("not src %s and not dst %s", intIpv4, intIpv4)); err != nil {
					panic(err)
				}

				packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

				for packet := range packetSource.Packets() {

					srcip := ""
					dstip := ""
					// fqdn := ""
					proto := ""
					sport := ""
					dport := ""

					if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
						ip, _ := ipLayer.(*layers.IPv4)
						srcip = ip.SrcIP.String()
						dstip = ip.DstIP.String()
					}

					if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
						proto = "TCP"
						tcp, _ := tcpLayer.(*layers.TCP)
						sport = tcp.SrcPort.String()
						dport = tcp.DstPort.String()
					}

					if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
						proto = "UDP"
						udp, _ := udpLayer.(*layers.UDP)
						sport = udp.SrcPort.String()
						dport = udp.DstPort.String()
					}

					netLogModel := &model.NetworkLog{
						ID:                 uuid.New().String(),
						TransactionAt:      packet.Metadata().Timestamp.In(vars.TZ),
						Protocol:           proto,
						SourceNetwork:      srcip,
						SourcePort:         sport,
						DestinationNetwork: dstip,
						DestinationPort:    dport,
					}

					if strings.TrimSpace(netLogModel.SourceNetwork) != "" && strings.TrimSpace(netLogModel.DestinationNetwork) != "" {
						if !utils.IsPrivateIPAddress(srcip) {
							netLogModel.TrafficFromInternet = true
						} else {

							session.Instance.UpdateLastSeen(netLogModel.SourceNetwork)
							// sessionKey := fmt.Sprintf("session|*|%s", netLogModel.SourceNetwork)
							// mSession, err := utils.RedisGetInsideWildcard[inmemory_model.Session](context.Background(), vars.RedisCache, sessionKey)
							// if err == nil {
							// 	mSession.LastSeen = netLogModel.TransactionAt
							// 	utils.RedisSet(
							// 		context.Background(),
							// 		vars.RedisCache,
							// 		fmt.Sprintf("session|%s", mSession.ID),
							// 		mSession,
							// 		time.Duration(vars.Config.SessionIdle)*time.Minute,
							// 	)
							// }

						}

						if !vars.Config.DisabledNetworkCapture {
							vars.NetLogDatabase.Index(netLogModel.ID, netLogModel)
						}

					}

				}

			}

		}
	}(ctx)
}
