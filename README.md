# Features

- User Managagement
- Captive Portal
- Bandwidth Control
- Network Capture
- LDAP Integration
- Radius Integration
- Private DNS Resolver (CoreDNS)

# Installation guide
(Need root user) 

1. Change config environment is **install.sh** file


      `vi install.sh`


2. Run install script


      `./install.sh`

# Uninstallation guide
(Need root user) 

`./uninstall.sh`

# Custom templates

You can use templates from the templates directory to modify them. The relevant and necessary parameters are already in that file.
After that, copy your custom file to replace that file on server in directory **/usr/share/cocong/templates**
