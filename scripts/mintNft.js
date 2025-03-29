require("dotenv").config();
const { ethers } = require("ethers");

// Load environment variables
const PRIVATE_KEY = process.env.PRIVATE_KEY;
const RPC_URL = process.env.RPC_URL_SEPOLIA;
const CONTRACT_ADDRESS = process.env.CONTRACT_ADDRESS_SEPOLIA;
const fromAddress = process.env.MINT_TO;

// ABI (Minimal)
const ABI = [
  "function mint(address to)",
  "function balanceOf(address owner) view returns (uint256)",
  "function tokenOfOwnerByIndex(address owner, uint256 index) view returns (uint256)",
];

async function main() {
  if (!fromAddress) {
    console.error("❌ No MINT_TO address passed in env");
    process.exit(1);
  }
  
  // Set up provider and wallet
  const provider = new ethers.JsonRpcProvider(RPC_URL);
  const wallet = new ethers.Wallet(PRIVATE_KEY, provider);
  
  console.log(`Connected to wallet: ${wallet.address}`);
  console.log(`Minting to: ${fromAddress}`);

  // Connect to contract
  const contract = new ethers.Contract(CONTRACT_ADDRESS, ABI, wallet);

  // Check if wallet has already minted
  const balance = await contract.balanceOf(fromAddress)

  if (balance==0) {
    console.log(`Wallet ${fromAddress} has not minted`);
    console.log(`Minting...`);
    const tx = await contract.mint(fromAddress);
    await tx.wait();
    console.log(`Mint successful! Transaction hash: ${tx.hash}`);
    const tokenId = await contract.tokenOfOwnerByIndex(fromAddress, 0)
    console.log(`🎟️  Wallet ${fromAddress} has just minted token ID: ${tokenId.toString()}`);
  } else {
    const tokenId = await contract.tokenOfOwnerByIndex(fromAddress, 0)
    console.log(`🎟️  Wallet ${fromAddress} has already minted token ID: ${tokenId.toString()}`);
  }
}

// Run the script
main().catch(console.error);
