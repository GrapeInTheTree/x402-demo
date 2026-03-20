package quiz

// SolidityQuestions returns Solidity quiz questions tested via forge.
func SolidityQuestions() []Question {
	return []Question{
		solERC20Basic(),
		solERC20Allowance(),
		solEIP3009Transfer(),
		solPermit2Approve(),
	}
}

func solERC20Basic() Question {
	return Question{
		ID: "sol-erc20-basic", Title: "ERC-20: balanceOf & transfer",
		Difficulty: "easy", Category: "ERC-20", Language: LangSolidity,
		Description: `Implement a minimal ERC-20 token with:
- A constructor that mints initialSupply to the deployer
- balanceOf(address) → returns balance
- transfer(address to, uint256 amount) → moves tokens

This is the foundation of all token interactions in x402.`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract SimpleToken {
    mapping(address => uint256) private _balances;
    uint256 public totalSupply;

    constructor(uint256 initialSupply) {
        // TODO: Mint initialSupply to msg.sender
        // 1. Set _balances[msg.sender] = initialSupply
        // 2. Set totalSupply = initialSupply
    }

    function balanceOf(address account) public view returns (uint256) {
        // TODO: Return the balance of account
        return 0;
    }

    function transfer(address to, uint256 amount) public returns (bool) {
        // TODO: Transfer amount from msg.sender to 'to'
        // 1. Check msg.sender has enough balance (require)
        // 2. Subtract from sender, add to receiver
        // 3. Return true
        return false;
    }
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract SimpleTokenTest is Test {
    SimpleToken token;
    address alice = address(1);
    address bob = address(2);

    function setUp() public {
        vm.prank(alice);
        token = new SimpleToken(1000000);
    }

    function test_InitialBalance() public view {
        assertEq(token.balanceOf(alice), 1000000);
        assertEq(token.balanceOf(bob), 0);
    }

    function test_TotalSupply() public view {
        assertEq(token.totalSupply(), 1000000);
    }

    function test_Transfer() public {
        vm.prank(alice);
        bool ok = token.transfer(bob, 100000);
        assertTrue(ok);
        assertEq(token.balanceOf(alice), 900000);
        assertEq(token.balanceOf(bob), 100000);
    }

    function test_TransferInsufficientBalance() public {
        vm.prank(bob);
        vm.expectRevert();
        token.transfer(alice, 1);
    }
}
`,
		Hints: []string{
			"constructor: _balances[msg.sender] = initialSupply; totalSupply = initialSupply;",
			"balanceOf: return _balances[account];",
			"transfer: require(_balances[msg.sender] >= amount); then subtract and add",
		},
	}
}

func solERC20Allowance() Question {
	return Question{
		ID: "sol-erc20-allowance", Title: "ERC-20: approve & transferFrom",
		Difficulty: "medium", Category: "ERC-20", Language: LangSolidity,
		Description: `Extend the ERC-20 with approval mechanics:
- approve(spender, amount) → allows spender to spend your tokens
- allowance(owner, spender) → returns current allowance
- transferFrom(from, to, amount) → spender moves tokens on behalf of owner

This is how Permit2 and EIP-3009 work under the hood.`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract TokenWithApproval {
    mapping(address => uint256) private _balances;
    mapping(address => mapping(address => uint256)) private _allowances;
    uint256 public totalSupply;

    constructor(uint256 initialSupply) {
        _balances[msg.sender] = initialSupply;
        totalSupply = initialSupply;
    }

    function balanceOf(address account) public view returns (uint256) {
        return _balances[account];
    }

    function transfer(address to, uint256 amount) public returns (bool) {
        require(_balances[msg.sender] >= amount, "insufficient balance");
        _balances[msg.sender] -= amount;
        _balances[to] += amount;
        return true;
    }

    function approve(address spender, uint256 amount) public returns (bool) {
        // TODO: Set allowance for spender to spend msg.sender's tokens
        // Store in _allowances[msg.sender][spender]
        return false;
    }

    function allowance(address owner, address spender) public view returns (uint256) {
        // TODO: Return the allowance
        return 0;
    }

    function transferFrom(address from, address to, uint256 amount) public returns (bool) {
        // TODO: Transfer tokens from 'from' to 'to' on behalf of msg.sender
        // 1. Check allowance: _allowances[from][msg.sender] >= amount
        // 2. Check balance: _balances[from] >= amount
        // 3. Subtract allowance, subtract balance, add to receiver
        return false;
    }
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract TokenApprovalTest is Test {
    TokenWithApproval token;
    address alice = address(1);
    address bob = address(2);
    address charlie = address(3);

    function setUp() public {
        vm.prank(alice);
        token = new TokenWithApproval(1000000);
    }

    function test_Approve() public {
        vm.prank(alice);
        assertTrue(token.approve(bob, 500000));
        assertEq(token.allowance(alice, bob), 500000);
    }

    function test_TransferFrom() public {
        vm.prank(alice);
        token.approve(bob, 200000);

        vm.prank(bob);
        assertTrue(token.transferFrom(alice, charlie, 100000));

        assertEq(token.balanceOf(alice), 900000);
        assertEq(token.balanceOf(charlie), 100000);
        assertEq(token.allowance(alice, bob), 100000);
    }

    function test_TransferFrom_ExceedAllowance() public {
        vm.prank(alice);
        token.approve(bob, 100);

        vm.prank(bob);
        vm.expectRevert();
        token.transferFrom(alice, charlie, 200);
    }

    function test_TransferFrom_InsufficientBalance() public {
        vm.prank(alice);
        token.approve(bob, type(uint256).max);

        vm.prank(bob);
        vm.expectRevert();
        token.transferFrom(alice, charlie, 2000000);
    }
}
`,
		Hints: []string{
			"approve: _allowances[msg.sender][spender] = amount; return true;",
			"allowance: return _allowances[owner][spender];",
			"transferFrom: check allowance, check balance, then _allowances[from][msg.sender] -= amount",
		},
	}
}

func solEIP3009Transfer() Question {
	return Question{
		ID: "sol-eip3009", Title: "EIP-3009: transferWithAuthorization Interface",
		Difficulty: "hard", Category: "EIP-3009", Language: LangSolidity,
		Description: `EIP-3009 adds transferWithAuthorization to USDC, enabling
gasless transfers via off-chain signatures.

Implement the authorization storage and validation logic.
(Signature verification is simplified — focus on the nonce
tracking and parameter validation.)`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract EIP3009Token {
    mapping(address => uint256) private _balances;
    mapping(address => mapping(bytes32 => bool)) private _usedNonces;
    uint256 public totalSupply;

    constructor(uint256 initialSupply) {
        _balances[msg.sender] = initialSupply;
        totalSupply = initialSupply;
    }

    function balanceOf(address account) public view returns (uint256) {
        return _balances[account];
    }

    // Simplified transferWithAuthorization (no actual signature verification)
    // In real USDC, this verifies an EIP-712 signature.
    function transferWithAuthorization(
        address from,
        address to,
        uint256 value,
        uint256 validAfter,
        uint256 validBefore,
        bytes32 nonce
    ) external returns (bool) {
        // TODO: 1. Check block.timestamp > validAfter
        // TODO: 2. Check block.timestamp < validBefore
        // TODO: 3. Check nonce not already used: !_usedNonces[from][nonce]
        // TODO: 4. Mark nonce as used
        // TODO: 5. Check from has sufficient balance
        // TODO: 6. Transfer: subtract from, add to
        // TODO: 7. Return true
        return false;
    }

    // Check if a nonce has been used
    function isNonceUsed(address account, bytes32 nonce) public view returns (bool) {
        // TODO: Return whether this nonce was already used
        return false;
    }
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

contract EIP3009Test is Test {
    EIP3009Token token;
    address alice = address(1);
    address bob = address(2);
    address facilitator = address(3);

    function setUp() public {
        vm.prank(alice);
        token = new EIP3009Token(1000000);
    }

    function test_TransferWithAuth() public {
        vm.warp(100);
        vm.prank(facilitator);
        bool ok = token.transferWithAuthorization(
            alice, bob, 100000,
            0,         // validAfter
            200,       // validBefore
            bytes32(uint256(1))  // nonce
        );
        assertTrue(ok);
        assertEq(token.balanceOf(alice), 900000);
        assertEq(token.balanceOf(bob), 100000);
    }

    function test_NonceReuse() public {
        bytes32 nonce = bytes32(uint256(42));
        vm.warp(100);

        vm.prank(facilitator);
        token.transferWithAuthorization(alice, bob, 1000, 0, 200, nonce);

        vm.prank(facilitator);
        vm.expectRevert();
        token.transferWithAuthorization(alice, bob, 1000, 0, 200, nonce);
    }

    function test_Expired() public {
        vm.warp(300);
        vm.prank(facilitator);
        vm.expectRevert();
        token.transferWithAuthorization(alice, bob, 1000, 0, 200, bytes32(uint256(2)));
    }

    function test_NotYetValid() public {
        vm.warp(50);
        vm.prank(facilitator);
        vm.expectRevert();
        token.transferWithAuthorization(alice, bob, 1000, 100, 200, bytes32(uint256(3)));
    }

    function test_IsNonceUsed() public {
        bytes32 nonce = bytes32(uint256(99));
        assertFalse(token.isNonceUsed(alice, nonce));

        vm.warp(100);
        vm.prank(facilitator);
        token.transferWithAuthorization(alice, bob, 100, 0, 200, nonce);

        assertTrue(token.isNonceUsed(alice, nonce));
    }
}
`,
		Hints: []string{
			"require(block.timestamp > validAfter, \"not yet valid\");",
			"require(!_usedNonces[from][nonce], \"nonce used\"); then _usedNonces[from][nonce] = true;",
			"require(_balances[from] >= value); _balances[from] -= value; _balances[to] += value;",
		},
	}
}

func solPermit2Approve() Question {
	return Question{
		ID: "sol-permit2-flow", Title: "Permit2: Token Approval Pattern",
		Difficulty: "medium", Category: "Permit2", Language: LangSolidity,
		Description: `Before using Permit2, the token owner must approve the Permit2
contract to spend their tokens. Then Permit2 can transfer tokens
on behalf of the owner to any spender.

Implement a mock Permit2 that:
1. Checks the token allowance from owner to Permit2 (this contract)
2. Calls transferFrom on the token to move funds`,
		Template: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

interface IERC20 {
    function transferFrom(address from, address to, uint256 amount) external returns (bool);
    function allowance(address owner, address spender) external view returns (uint256);
}

contract MockPermit2 {
    // Transfer tokens from owner to recipient via Permit2.
    // This contract must have been approved by the owner on the token contract.
    function transferFrom(
        address token,
        address owner,
        address recipient,
        uint256 amount
    ) external returns (bool) {
        // TODO: 1. Check that owner has approved this contract (Permit2) for >= amount
        //          Use IERC20(token).allowance(owner, address(this))
        // TODO: 2. Call IERC20(token).transferFrom(owner, recipient, amount)
        // TODO: 3. Return true
        return false;
    }

    // Check if owner has sufficient allowance for Permit2
    function hasAllowance(
        address token,
        address owner,
        uint256 amount
    ) external view returns (bool) {
        // TODO: Return whether allowance >= amount
        return false;
    }
}
`,
		TestCode: `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/Solution.sol";

// Simple ERC20 for testing
contract TestToken {
    mapping(address => uint256) public balanceOf;
    mapping(address => mapping(address => uint256)) public allowance;

    constructor(address to, uint256 amount) {
        balanceOf[to] = amount;
    }

    function approve(address spender, uint256 amount) external returns (bool) {
        allowance[msg.sender][spender] = amount;
        return true;
    }

    function transferFrom(address from, address to, uint256 amount) external returns (bool) {
        require(allowance[from][msg.sender] >= amount, "allowance");
        require(balanceOf[from] >= amount, "balance");
        allowance[from][msg.sender] -= amount;
        balanceOf[from] -= amount;
        balanceOf[to] += amount;
        return true;
    }
}

contract MockPermit2Test is Test {
    MockPermit2 permit2;
    TestToken token;
    address alice = address(1);
    address bob = address(2);

    function setUp() public {
        permit2 = new MockPermit2();
        token = new TestToken(alice, 1000000);

        // Alice approves Permit2 contract
        vm.prank(alice);
        token.approve(address(permit2), 500000);
    }

    function test_HasAllowance() public view {
        assertTrue(permit2.hasAllowance(address(token), alice, 100000));
        assertFalse(permit2.hasAllowance(address(token), alice, 600000));
    }

    function test_TransferFrom() public {
        bool ok = permit2.transferFrom(address(token), alice, bob, 100000);
        assertTrue(ok);
        assertEq(token.balanceOf(bob), 100000);
        assertEq(token.balanceOf(alice), 900000);
    }

    function test_TransferFrom_NoAllowance() public {
        vm.expectRevert();
        permit2.transferFrom(address(token), bob, alice, 100);
    }
}
`,
		Hints: []string{
			"IERC20(token).allowance(owner, address(this)) gives the Permit2 allowance",
			"require(IERC20(token).allowance(owner, address(this)) >= amount);",
			"return IERC20(token).transferFrom(owner, recipient, amount);",
		},
	}
}
