# gofind

```
https://www.linuxbabe.com/linux-server/set-up-package-repository-debian-ubuntu-server && \
sudo echo "deb [arch=amd64] http://marten.zip/gf jammy main" > /etc/apt/sources.list.d/gf.list && \
curl -sS https://marten.zip/gf/gpg-pubkey.asc | gpg --dearmor | sudo tee /etc/apt/trusted.gpg.d/marten.gpg && \
sudo apt update && \
sudo apt install gf
```
