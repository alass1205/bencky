#!/bin/sh

# Vérifier si le nœud est déjà initialisé
if [ ! -f /root/.ethereum/geth/LOCK ]; then
    echo "Initializing Geth with genesis block..."
    geth init /root/genesis.json --datadir /root/.ethereum
fi

# Exécuter Geth avec les paramètres passés
exec geth "$@"
