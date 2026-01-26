// Binary Tree Operations Benchmark for goTS
// A comprehensive BST implementation with various tree operations

// TreeNode class - uses sentinel pattern with value -1 to indicate "empty"
// This avoids nullable type access issues
class TreeNode {
    value: int
    hasLeft: boolean
    hasRight: boolean
    leftNode: TreeNode | null
    rightNode: TreeNode | null

    constructor(value: int) {
        this.value = value
        this.hasLeft = false
        this.hasRight = false
        this.leftNode = null
        this.rightNode = null
    }

    // Set left child
    setLeft(child: TreeNode): void {
        this.leftNode = child
        this.hasLeft = true
    }

    // Set right child
    setRight(child: TreeNode): void {
        this.rightNode = child
        this.hasRight = true
    }

    // Get left child (returns self as sentinel if no left child)
    left(): TreeNode {
        let child: TreeNode | null = this.leftNode
        if (this.hasLeft && child != null) {
            return child
        }
        return this
    }

    // Get right child (returns self as sentinel if no right child)
    right(): TreeNode {
        let child: TreeNode | null = this.rightNode
        if (this.hasRight && child != null) {
            return child
        }
        return this
    }

    // Check if leaf node
    isLeaf(): boolean {
        return !this.hasLeft && !this.hasRight
    }
}

// Binary Search Tree class
class BST {
    rootNode: TreeNode | null
    nodeCount: int

    constructor() {
        this.rootNode = null
        this.nodeCount = 0
    }

    // Check if tree is empty
    isEmpty(): boolean {
        return this.rootNode == null
    }

    // Get root (creates sentinel if null)
    root(): TreeNode {
        let r: TreeNode | null = this.rootNode
        if (r != null) {
            return r
        }
        // This shouldn't be called on empty tree
        return new TreeNode(-1)
    }

    // Insert a value
    insert(value: int): void {
        if (this.rootNode == null) {
            this.rootNode = new TreeNode(value)
            this.nodeCount = 1
        } else {
            this.insertAt(this.root(), value)
            this.nodeCount = this.nodeCount + 1
        }
    }

    insertAt(node: TreeNode, value: int): void {
        if (value < node.value) {
            if (node.hasLeft) {
                this.insertAt(node.left(), value)
            } else {
                node.setLeft(new TreeNode(value))
            }
        } else if (value > node.value) {
            if (node.hasRight) {
                this.insertAt(node.right(), value)
            } else {
                node.setRight(new TreeNode(value))
            }
        }
        // Equal values are ignored (no duplicates)
    }

    // Search for a value
    search(value: int): boolean {
        if (this.isEmpty()) {
            return false
        }
        return this.searchAt(this.root(), value)
    }

    searchAt(node: TreeNode, value: int): boolean {
        if (value == node.value) {
            return true
        }
        if (value < node.value && node.hasLeft) {
            return this.searchAt(node.left(), value)
        }
        if (value > node.value && node.hasRight) {
            return this.searchAt(node.right(), value)
        }
        return false
    }

    // Find minimum value
    findMin(): int {
        if (this.isEmpty()) {
            return -1
        }
        return this.findMinAt(this.root())
    }

    findMinAt(node: TreeNode): int {
        if (node.hasLeft) {
            return this.findMinAt(node.left())
        }
        return node.value
    }

    // Find maximum value
    findMax(): int {
        if (this.isEmpty()) {
            return -1
        }
        return this.findMaxAt(this.root())
    }

    findMaxAt(node: TreeNode): int {
        if (node.hasRight) {
            return this.findMaxAt(node.right())
        }
        return node.value
    }

    // Calculate height
    height(): int {
        if (this.isEmpty()) {
            return 0
        }
        return this.heightAt(this.root())
    }

    heightAt(node: TreeNode): int {
        let leftH: int = 0
        let rightH: int = 0
        if (node.hasLeft) {
            leftH = this.heightAt(node.left())
        }
        if (node.hasRight) {
            rightH = this.heightAt(node.right())
        }
        if (leftH > rightH) {
            return 1 + leftH
        }
        return 1 + rightH
    }

    // Count all nodes
    countNodes(): int {
        if (this.isEmpty()) {
            return 0
        }
        return this.countAt(this.root())
    }

    countAt(node: TreeNode): int {
        let count: int = 1
        if (node.hasLeft) {
            count = count + this.countAt(node.left())
        }
        if (node.hasRight) {
            count = count + this.countAt(node.right())
        }
        return count
    }

    // Count leaf nodes
    countLeaves(): int {
        if (this.isEmpty()) {
            return 0
        }
        return this.countLeavesAt(this.root())
    }

    countLeavesAt(node: TreeNode): int {
        if (node.isLeaf()) {
            return 1
        }
        let count: int = 0
        if (node.hasLeft) {
            count = count + this.countLeavesAt(node.left())
        }
        if (node.hasRight) {
            count = count + this.countLeavesAt(node.right())
        }
        return count
    }

    // Sum all values
    sum(): int {
        if (this.isEmpty()) {
            return 0
        }
        return this.sumAt(this.root())
    }

    sumAt(node: TreeNode): int {
        let total: int = node.value
        if (node.hasLeft) {
            total = total + this.sumAt(node.left())
        }
        if (node.hasRight) {
            total = total + this.sumAt(node.right())
        }
        return total
    }

    // Check if balanced
    isBalanced(): boolean {
        if (this.isEmpty()) {
            return true
        }
        return this.checkBalanceAt(this.root()) != -1
    }

    checkBalanceAt(node: TreeNode): int {
        let leftH: int = 0
        let rightH: int = 0
        if (node.hasLeft) {
            leftH = this.checkBalanceAt(node.left())
            if (leftH == -1) {
                return -1
            }
        }
        if (node.hasRight) {
            rightH = this.checkBalanceAt(node.right())
            if (rightH == -1) {
                return -1
            }
        }
        let diff: int = leftH - rightH
        if (diff < 0) {
            diff = 0 - diff
        }
        if (diff > 1) {
            return -1
        }
        if (leftH > rightH) {
            return 1 + leftH
        }
        return 1 + rightH
    }

    // Validate BST property
    isValidBST(): boolean {
        if (this.isEmpty()) {
            return true
        }
        return this.validateAt(this.root(), -2147483647, 2147483647)
    }

    validateAt(node: TreeNode, minVal: int, maxVal: int): boolean {
        if (node.value <= minVal || node.value >= maxVal) {
            return false
        }
        let leftValid: boolean = true
        let rightValid: boolean = true
        if (node.hasLeft) {
            leftValid = this.validateAt(node.left(), minVal, node.value)
        }
        if (node.hasRight) {
            rightValid = this.validateAt(node.right(), node.value, maxVal)
        }
        return leftValid && rightValid
    }

    // Inorder traversal
    inorderTraversal(): int[] {
        let result: int[] = []
        if (!this.isEmpty()) {
            this.collectInorder(this.root(), result)
        }
        return result
    }

    collectInorder(node: TreeNode, result: int[]): void {
        if (node.hasLeft) {
            this.collectInorder(node.left(), result)
        }
        result.push(node.value)
        if (node.hasRight) {
            this.collectInorder(node.right(), result)
        }
    }

    // Preorder traversal
    preorderTraversal(): int[] {
        let result: int[] = []
        if (!this.isEmpty()) {
            this.collectPreorder(this.root(), result)
        }
        return result
    }

    collectPreorder(node: TreeNode, result: int[]): void {
        result.push(node.value)
        if (node.hasLeft) {
            this.collectPreorder(node.left(), result)
        }
        if (node.hasRight) {
            this.collectPreorder(node.right(), result)
        }
    }

    // Postorder traversal
    postorderTraversal(): int[] {
        let result: int[] = []
        if (!this.isEmpty()) {
            this.collectPostorder(this.root(), result)
        }
        return result
    }

    collectPostorder(node: TreeNode, result: int[]): void {
        if (node.hasLeft) {
            this.collectPostorder(node.left(), result)
        }
        if (node.hasRight) {
            this.collectPostorder(node.right(), result)
        }
        result.push(node.value)
    }

    // Find LCA (Lowest Common Ancestor)
    lca(a: int, b: int): int {
        if (this.isEmpty()) {
            return -1
        }
        return this.findLCA(this.root(), a, b)
    }

    findLCA(node: TreeNode, a: int, b: int): int {
        if (a < node.value && b < node.value && node.hasLeft) {
            return this.findLCA(node.left(), a, b)
        }
        if (a > node.value && b > node.value && node.hasRight) {
            return this.findLCA(node.right(), a, b)
        }
        return node.value
    }

    // Check if path sum exists
    hasPathSum(targetSum: int): boolean {
        if (this.isEmpty()) {
            return false
        }
        return this.checkPathSum(this.root(), targetSum)
    }

    checkPathSum(node: TreeNode, remaining: int): boolean {
        let newRemaining: int = remaining - node.value
        if (node.isLeaf()) {
            return newRemaining == 0
        }
        let leftHas: boolean = false
        let rightHas: boolean = false
        if (node.hasLeft) {
            leftHas = this.checkPathSum(node.left(), newRemaining)
        }
        if (node.hasRight) {
            rightHas = this.checkPathSum(node.right(), newRemaining)
        }
        return leftHas || rightHas
    }

    // Get nodes at specific level
    getLevel(level: int): int[] {
        let result: int[] = []
        if (!this.isEmpty()) {
            this.collectLevel(this.root(), level, 0, result)
        }
        return result
    }

    collectLevel(node: TreeNode, target: int, current: int, result: int[]): void {
        if (current == target) {
            result.push(node.value)
        } else {
            if (node.hasLeft) {
                this.collectLevel(node.left(), target, current + 1, result)
            }
            if (node.hasRight) {
                this.collectLevel(node.right(), target, current + 1, result)
            }
        }
    }

    // Level widths
    levelWidths(): int[] {
        let widths: int[] = []
        let h: int = this.height()
        let i: int = 0
        while (i < h) {
            let level: int[] = this.getLevel(i)
            widths.push(level.length)
            i = i + 1
        }
        return widths
    }

    // Delete a value
    delete(value: int): boolean {
        if (this.isEmpty()) {
            return false
        }
        if (!this.search(value)) {
            return false
        }
        this.rootNode = this.deleteAt(this.root(), value)
        this.nodeCount = this.nodeCount - 1
        return true
    }

    deleteAt(node: TreeNode, value: int): TreeNode | null {
        if (value < node.value && node.hasLeft) {
            let newLeft: TreeNode | null = this.deleteAt(node.left(), value)
            if (newLeft != null) {
                node.setLeft(newLeft)
            } else {
                node.hasLeft = false
                node.leftNode = null
            }
        } else if (value > node.value && node.hasRight) {
            let newRight: TreeNode | null = this.deleteAt(node.right(), value)
            if (newRight != null) {
                node.setRight(newRight)
            } else {
                node.hasRight = false
                node.rightNode = null
            }
        } else if (value == node.value) {
            // Node to delete found
            if (!node.hasLeft && !node.hasRight) {
                return null
            }
            if (!node.hasLeft) {
                return node.rightNode
            }
            if (!node.hasRight) {
                return node.leftNode
            }
            // Two children - find inorder successor
            let successorVal: int = this.findMinAt(node.right())
            node.value = successorVal
            let newRight: TreeNode | null = this.deleteAt(node.right(), successorVal)
            if (newRight != null) {
                node.setRight(newRight)
            } else {
                node.hasRight = false
                node.rightNode = null
            }
        }
        return node
    }

    // Calculate diameter
    diameter(): int {
        if (this.isEmpty()) {
            return 0
        }
        let maxDiam: int = 0
        this.calcDiameter(this.root(), maxDiam)
        return maxDiam
    }

    calcDiameter(node: TreeNode, maxD: int): int {
        let leftH: int = 0
        let rightH: int = 0
        if (node.hasLeft) {
            leftH = this.calcDiameter(node.left(), maxD)
        }
        if (node.hasRight) {
            rightH = this.calcDiameter(node.right(), maxD)
        }
        let currDiam: int = leftH + rightH
        if (currDiam > maxD) {
            maxD = currDiam
        }
        if (leftH > rightH) {
            return 1 + leftH
        }
        return 1 + rightH
    }
}

// Linear Congruential Generator for pseudo-random numbers
class RNG {
    state: int

    constructor(seed: int) {
        this.state = seed
    }

    next(): int {
        // LCG parameters (glibc)
        let a: int = 1103515245
        let c: int = 12345
        let m: int = 2147483647
        this.state = (a * this.state + c) % m
        if (this.state < 0) {
            this.state = 0 - this.state
        }
        return this.state
    }

    nextInRange(maxVal: int): int {
        let v: int = this.next() % maxVal
        if (v < 0) {
            v = 0 - v
        }
        return v
    }
}

// Benchmark helper - convert bool to string
function boolStr(val: boolean): string {
    if (val) {
        return "true"
    }
    return "false"
}

// Main benchmark function
function benchmark(n: int): void {
    console.log("=== Binary Tree Benchmark (" + (n.toString()) + " elements) ===")
    console.log("")

    let tree: BST = new BST()
    let rng: RNG = new RNG(42)

    // 1. Insert elements
    console.log("1. Inserting elements...")
    let i: int = 0
    while (i < n) {
        let v: int = rng.nextInRange(100000)
        tree.insert(v)
        i = i + 1
    }
    console.log("   Tree size: " + (tree.countNodes()).toString())

    // 2. Tree properties
    console.log("")
    console.log("2. Tree Properties:")
    console.log("   Height: " + (tree.height()).toString())
    console.log("   Leaves: " + (tree.countLeaves()).toString())
    console.log("   Balanced: " + boolStr(tree.isBalanced()))
    console.log("   Valid BST: " + boolStr(tree.isValidBST()))
    console.log("   Sum: " + (tree.sum()).toString())
    console.log("   Min: " + (tree.findMin()).toString())
    console.log("   Max: " + (tree.findMax()).toString())

    // 3. Level widths
    console.log("")
    console.log("3. Level widths (first 10):")
    let widths: int[] = tree.levelWidths()
    let j: int = 0
    while (j < 10 && j < widths.length) {
        console.log("   Level " + (j.toString()) + ": " + (widths[j].toString()))
        j = j + 1
    }

    // 4. Search operations
    console.log("")
    console.log("4. Searching 100 random values...")
    let found: int = 0
    let searchRng: RNG = new RNG(12345)
    let k: int = 0
    while (k < 100) {
        let sv: int = searchRng.nextInRange(100000)
        if (tree.search(sv)) {
            found = found + 1
        }
        k = k + 1
    }
    console.log("   Found: " + (found.toString()) + "/100")

    // 5. LCA
    console.log("")
    console.log("5. LCA of min and max:")
    let minV: int = tree.findMin()
    let maxV: int = tree.findMax()
    let lcaV: int = tree.lca(minV, maxV)
    console.log("   LCA(" + (minV.toString()) + "," + (maxV.toString()) + ") = " + (lcaV.toString()))

    // 6. Path sum
    console.log("")
    console.log("6. Path sum check:")
    let nc: int = tree.countNodes()
    let avgH: int = tree.height()
    let testS: int = parseInt(tree.sum() / nc) * avgH
    console.log("   Has path sum " + (testS.toString()) + ": " + boolStr(tree.hasPathSum(testS)))

    // 7. Traversals
    console.log("")
    console.log("7. Traversals (first 5 elements):")
    let inord: int[] = tree.inorderTraversal()
    let preord: int[] = tree.preorderTraversal()
    let postord: int[] = tree.postorderTraversal()

    console.log("   Inorder:   ")
    let m: int = 0
    while (m < 5 && m < inord.length) {
        console.log((inord[m].toString()) + " ")
        m = m + 1
    }
    console.log("")

    console.log("   Preorder:  ")
    m = 0
    while (m < 5 && m < preord.length) {
        console.log((preord[m].toString()) + " ")
        m = m + 1
    }
    console.log("")

    console.log("   Postorder: ")
    let plen: int = postord.length
    let si: int = plen - 5
    if (si < 0) {
        si = 0
    }
    m = si
    while (m < plen) {
        console.log((postord[m].toString()) + " ")
        m = m + 1
    }
    console.log("")

    // 8. Verify sorted
    console.log("")
    console.log("8. Verify BST property:")
    let sorted: boolean = true
    let s: int = 1
    while (s < inord.length) {
        if (inord[s] < inord[s - 1]) {
            sorted = false
        }
        s = s + 1
    }
    console.log("   Inorder sorted: " + boolStr(sorted))

    // 9. Delete operations
    console.log("")
    console.log("9. Deleting 50 random elements...")
    let deleted: int = 0
    let delRng: RNG = new RNG(99999)
    let d: int = 0
    while (d < 50) {
        let dv: int = delRng.nextInRange(100000)
        if (tree.delete(dv)) {
            deleted = deleted + 1
        }
        d = d + 1
    }
    console.log("   Deleted: " + (deleted.toString()))
    console.log("   New size: " + (tree.countNodes().toString()))
    console.log("   Still valid: " + boolStr(tree.isValidBST()))

    // 10. Additional metrics
    console.log("")
    console.log("10. Post-deletion metrics:")
    console.log("   Height: " + (tree.height().toString()))
    console.log("   Leaves: " + (tree.countLeaves().toString()))

    let level5: int[] = tree.getLevel(5)
    console.log("   Nodes at level 5: " + (level5.length.toString()))

    console.log("")
    console.log("=== Benchmark Complete ===")
}

// Run benchmarks
console.log("Binary Tree Operations Benchmark")
console.log("================================")
console.log("")

benchmark(1000000)
