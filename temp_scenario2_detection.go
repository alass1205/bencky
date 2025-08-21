func (nm *NetworkMonitor) hasScenario2BeenExecuted() bool {
	aliceTxCount := nm.getTransactionCount("http://localhost:8545", "0x71562b71999873db5b286df957af199ec94617f7")
	return aliceTxCount >= 3
}
