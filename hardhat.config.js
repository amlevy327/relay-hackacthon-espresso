require("@nomicfoundation/hardhat-toolbox");
require("dotenv").config();

module.exports = {
  solidity: "0.8.28",
  networks: {
    rollupA: {
      url: process.env.RPC_URL,
      accounts: [process.env.PRIVATE_KEY]
    },
    sepolia: {
      url: process.env.RPC_URL_SEPOLIA,
      accounts: [process.env.PRIVATE_KEY]
    },
  }
};