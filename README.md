This project is an extension of: https://github.com/EspressoSystems/hackathon-example

🛡️ Rollup Transaction Monitor with NFT Minting

Monitors transactions on a custom Arbitrum Sepolia Nitro rollup (ID: 327327327), filtering for transactions of at least 1 wei sent to a specific recipient. When a matching transaction is detected, the relay mints an NFT on Ethereum Sepolia to the sender.

🚀 Try It Out

Live Demo: https://espresso-stellar-tau.vercel.app/

🔍 What It Does

✅ Block Monitoring
Retrieves the latest processed block from a Caff node and inspects transactions.

✅ Transaction Filtering
Filters based on transaction 'value' and 'to' address.

✅ NFT Minting
Triggers an NFT mint on Ethereum Sepolia to transaction 'from' address when meeting criteria.