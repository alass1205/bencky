# 📋 Benchy Audit Verification Guide

Ce guide fournit les commandes exactes pour vérifier chaque critère d'audit du projet Benchy.

## 🔧 **Critères Techniques**

### ✅ **1. Are two different clients launched?**

**Vérification :**
```bash
# Vérifier les images Docker utilisées
cat docker/docker-compose.yml | grep "image:"

# Vérifier les clients en cours d'exécution
docker ps --format "table {{.Names}}\t{{.Image}}" | grep benchy
```

**Résultat attendu :**
- Geth : benchy-alice, benchy-cassandra, benchy-elena
- Nethermind : benchy-bob, benchy-driss

---

### ✅ **2. Is clique the consensus algorithm?**

**Vérification :**
```bash
# Vérifier Network ID commun (preuve de réseau partagé)
echo "🔍 Network ID verification:"
for port in 8545 8547 8549 8551 8553; do
  result=$(curl -s -X POST -H "Content-Type: application/json" \
    --data '{"jsonrpc":"2.0","method":"net_version","params":[],"id":1}' \
    http://localhost:$port | jq -r '.result // "offline"')
  echo "Port $port: Network ID $result"
done

# Vérifier l'affichage Clique dans le monitoring
./bin/benchy infos | grep "Consensus: Clique PoA"
```

**Résultat attendu :**
- Tous les nœuds : Network ID `1337`
- Affichage : `🔗 Consensus: Clique PoA | Network ID: 12345 | Validators: Alice, Bob, Cassandra`

---

## 🚀 **Launch Network & Status**

### ✅ **3. Launch the network and display status**

**Vérification :**
```bash
# Nettoyer et lancer
docker-compose -f docker/docker-compose.yml down -v
make build
./bin/benchy launch-network

# Afficher le statut
./bin/benchy infos
```

---

### ✅ **4. Does the command launch the five nodes?**

**Vérification :**
```bash
# Compter les conteneurs lancés
docker ps | grep benchy | wc -l

# Lister les nœuds
docker ps --format "table {{.Names}}\t{{.Status}}" | grep benchy
```

**Résultat attendu :** 5 conteneurs actifs

---

### ✅ **5. Does the interface display the latest block of each node?**

**Vérification :**
```bash
# Vérifier colonne Block dans infos
./bin/benchy infos | grep -E "Block|#[0-9]+"
```

**Résultat attendu :** Colonne "Block" avec numéros (ex: #0, #3, #6)

---

### ✅ **6. Does the interface display their Ethereum address and balance?**

**Vérification :**
```bash
# Vérifier colonne Balance dans infos
./bin/benchy infos | grep -E "Balance|ETH"
```

**Résultat attendu :** Balances ETH affichées (ex: 100.0000 ETH, 99.7000 ETH)

---

### ✅ **7. Does the interface display CPU and memory consumption?**

**Vérification :**
```bash
# Vérifier colonnes CPU% et Memory
./bin/benchy infos | grep -E "CPU%|Memory|MiB"
```

**Résultat attendu :** Colonnes CPU% et Memory avec valeurs (ex: 0.1%, 32.09MiB)

---

## 🎬 **Scenarios Testing**

### ✅ **8. Scenario 0 - Alice, Bob, Cassandra positive balance?**

**Vérification :**
```bash
./bin/benchy scenario 0
./bin/benchy infos | grep -E "Alice|Bob|Cassandra" | grep ETH
```

**Résultat attendu :** Balances positives pour les 3 validateurs

---

### ✅ **9. Scenario 1 - Feedback and updated balances?**

**Vérification :**
```bash
# Sauvegarder état initial
./bin/benchy infos > before_scenario1.txt

# Exécuter scénario 1
./bin/benchy scenario 1

# Vérifier feedback (doit afficher progression des transactions)
# Comparer balances
./bin/benchy infos > after_scenario1.txt
diff before_scenario1.txt after_scenario1.txt
```

**Résultat attendu :**
- Feedback détaillé des transactions
- Alice : balance diminuée (-0.3 ETH)
- Bob : balance augmentée (+0.3 ETH)

---

### ✅ **10. Mempool displayed?**

**Vérification :**
```bash
# Vérifier colonne Mempool
./bin/benchy infos | grep -E "Mempool|txs"
```

**Résultat attendu :** Colonne "Mempool" avec valeurs (ex: 1 txs, 2 txs)

---

### ✅ **11. Scenario 2 - Driss & Elena receive 1000 BY tokens?**

**Vérification :**
```bash
./bin/benchy scenario 2
./bin/benchy infos | grep -E "Driss|Elena" | grep "BY"
```

**Résultat attendu :** Driss et Elena affichent "1000 BY tokens"

---

### ✅ **12. Scenario 3 - Elena receives 1 ETH?**

**Vérification :**
```bash
# Sauvegarder balance Elena avant
./bin/benchy infos | grep Elena

./bin/benchy scenario 3

# Vérifier balance Elena après
./bin/benchy infos | grep Elena
```

**Résultat attendu :** Elena passe de "1000 BY tokens" à "1000 BY + 1.0 ETH"

---

## 🔧 **Temporary Failure Testing**

### ✅ **13-15. Temporary failure tests**

**Vérification :**
```bash
# État initial
./bin/benchy infos | grep Alice

# Déclencher panne (dans terminal 1)
./bin/benchy temporary-failure alice

# Pendant la panne (dans terminal 2)
./bin/benchy infos | grep Alice
# Résultat attendu: Alice 🔴 OFF

# Attendre 40 secondes, puis vérifier retour
./bin/benchy infos | grep Alice
# Résultat attendu: Alice 🟢 ON

# Vérifier mise à jour après 2 minutes
sleep 120
./bin/benchy infos | grep Alice
```

**Résultats attendus :**
- Pendant panne : `Alice 🔴 OFF`
- Après 40s : `Alice 🟢 ON`
- Bloc mis à jour quelques minutes plus tard

---

## 🎁 **Bonus**

### ✅ **16. Does `-u` option work for regular updates?**

**Vérification :**
```bash
# Test mise à jour continue (5 secondes)
./bin/benchy infos -u 5

# Test avec d'autres commandes
./bin/benchy infos --update 10
```

**Résultat attendu :** Affichage qui se met à jour automatiquement

---

## 📊 **Script d'audit complet**

```bash
#!/bin/bash
echo "🔍 BENCHY AUDIT AUTOMATION"
echo "=========================="

# 1. Setup
docker-compose -f docker/docker-compose.yml down -v
make build
./bin/benchy launch-network

# 2. Tests basiques
echo "✅ Network launched"
./bin/benchy infos

# 3. Tests scénarios
./bin/benchy scenario 0
./bin/benchy scenario 1
./bin/benchy scenario 2
./bin/benchy scenario 3

# 4. Test panne
echo "🔧 Testing temporary failure..."
./bin/benchy temporary-failure alice &
sleep 5
./bin/benchy infos | grep Alice

echo "🎯 Audit terminé - Vérifiez les résultats ci-dessus"
```

## 🎯 **Score attendu : 15/15 ✅**

Tous les critères sont vérifiables et fonctionnels dans le projet Benchy.