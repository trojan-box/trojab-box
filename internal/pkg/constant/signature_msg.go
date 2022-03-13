package constant

const LoginSignMsg = `
trojan-box wants you to sign in with your Ethereum account:
${address}
for login
Nonce: ${nonce}
Timestamp: ${timestamp}
`

const PlayGameSignMsg = `
trojan-box wants you to sign in with your Ethereum account:
${address}
for play game
Nonce: ${nonce}
Timestamp: ${timestamp}
Cards: ${cards}
Chosen: ${chosen}
`

const WithdrawBonusApplySignMsg = `
trojan-box wants you to sign in with your Ethereum account:
${address}
for withdraw bonus apply
Nonce: ${nonce}
Timestamp: ${timestamp}
Bonus: ${bonus}
`
