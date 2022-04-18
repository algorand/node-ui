package about

func GetHelpContent() string {
	return `
# Algorand Node UI :smiley_cat:

An **awesome** node Terminal User Interface for node runners.

Easy access to **important** tools and node information.

# Status

Continuous status is available for:
* Network information.
* Protocol upgrade status.
* Catchup sync time.
* Fast catchup progress.

# Explorer

## Blocks

Full real-time access to block information, and aggregations including:
* Number of transactions.
* Transaction types.
* Sum of payment transactions.	
* You get a gold star for actually reading this.
* Unique assets used in asset transactions.
* Unique applications used in applications.

## Transactions

Drill into a block for a detailed transaction breakdown:
* sender
* type
* transfer amount for payment / asset transfer transactions
* signature type, including inner-transactions

## Raw Transaction

View the raw transaction details.

# Utilities

Shortcuts for handy utilities.

# Accounts

View all of your accounts along with recent transactions.

# Configuration

Full node configuration details.
` + "```json" + `
{
    "Version": 16,
    "AccountsRebuildSynchronousMode": 1,
    "AnnounceParticipationKey": true,
    "Archival": false,
    "BaseLoggerDebugLevel": 4,
    "BroadcastConnectionsLimit": -1,
    "CadaverSizeTarget": 1073741824,
    ...
}
` + "```" + `

# Help

Let's be realistic for a moment, this software was so
intuitive and fun to use that you have no need for help.

But don't worry, it's right here if you need it!
`
}
