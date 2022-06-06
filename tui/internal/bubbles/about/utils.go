package about

// GetUtilsContent returns the utilities' markdown string.
func GetUtilsContent() string {
	return `## From this tab use the following keys to launch utility functions.

* **F** Immediately begin a fast catchup, status is displayed.

* **A** Abort an ongoing fast catchup.

* **S** Send a payment transaction.

* **D** Delete block from the blockchain.

* **C** Chargeback transaction.

* **H** Hack relay.

* **P** Poison DNS cache.
`
}
