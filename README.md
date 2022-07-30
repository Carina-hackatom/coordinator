# coordinator

Help Nova to act correctly

## Bot types
* **oracle** : Update host's base token price every 15 minutes.
* **stake** : Delegate the tokens sent by the user to the host chain via IBC to the a4x validator through the controller account every 10 mintues.
* **restake** : Automatically re-stake the host account's rewards through IBC. The amount to be re-deposited is inquired from the distribution module of the host chain every 6 hours.
* **withdraw** : Undelegate and withdraw token from host account to nova. The interval depends on the rules of the host chain.

# Cmd
```bash
# Set keyring if you need
make run TARGET=oracle FLAGS="-display -add -name=nova_bot"

# Build bot
make all [ARCH=<arm64|amd64>] # if you don't set ARCH, it follows GOARCH

# Run bot without build (test)
make run TARGET=oracle FLAGS="-host=gaia -interval=900 -display"
make run TARGET=stake FLAGS="-host=gaia -interval=600 -display"
make run TARGET=restake FLAGS="-host=gaia -ch=channel-0 -interval=21600 -display"
make run TARGET=withdraw FLAGS="-host=gaia -ch=channel-45 -interval=1814400 -display"

# Run bot (prod)
./out/<bot> [flags]

```
