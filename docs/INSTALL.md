----------
avorsmtp 
----------

 
### Installation 
Copier le ficheir avorsmtp_[version].tar.gz sur la machine cible dans le repertoire /opt/avorsmtp
Exécuter tar xvzf avorsmtp_[version].tar.gz

Ce placer dans le repertoire : cd /opt/avorsmtp/avorsmtp_[version]

> **NOTE:**
>
> - En cas d'une nouvelle installation copier le ficheir avorsmtp.supervisor.conf dans le répertoire /etc/supervisor/conf.d
>
> - Copier le fichier /opt/avorsmtp/avorsmtp_[version]/config.sample.json dans /opt/avorsmtp/avorsmtp_[version]/config.json
> - Addapter ce fichier en cas de besoin à votre environement cible


Vérifier si l'application est en cours d'execution via console de supervisor : supervisorctl
Si l'application est en court d'exécution arrêter l'application : stop avorsmtp
Quitter le console : exit

Crée un lien symbolic ln -s /opt/avorsmtp//avorsmtp_[version] /opt/avorsmtp/current 
> **NOTE:**
>
> - En cas si le répertoire existe /opt/avorsmtp/current. Supprimer rm -rf /opt/avorsmtp/current

### Configuraiton
Copier le fichier config.json de répértoire d'instalaltion dans le répértoire de l'application : cp ./samples/config.json config.json
Copier le fichier logger.xml de répértoire d'instalaltion dans le répértoire de l'application : cp ./samples/logger.xml logger.xml

Options de configuration de la connexion au serveur smtp

	"InsecureSkipVerify" : certificat TLS n'est pas verififé si la valeur de l'option est égal à true
	"UnencryptedAuth": force utiliser une connexion non sécuriser si la valeur de l'option est égal à true
	"SkipAuth": procedure d'authentification n'est pas utilisé si la valeur de l'option est égal à true
	
Paramètre templateFile permet d'utiliser un fichier text ou un fichier en format html pour la génération d'un e-mail.

Paramètre emailFormat definit le format d'e-mail à envoyer par defaut c'est text/plain. Pour envoie d'un e-mail en format html : text/html

### Note 
	En cas d'utilsiaiton d'un serveur smtp sur le port 465 merci de voir la configuraiton de relay sur le serveur postfix.


### Mise à jour 
Mise à jour est identique à l'installation sans la partie de la configuration.

Il faut copier le ficheir de la configuration actuel (/opt/avorsmtp/current/config.json) dans le repertoire /opt/avorsmtp/avorsmtp_[version] 





