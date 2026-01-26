package main

import (
	"fmt"
	"reflect"
)

// Runtime helpers

func gts_len(v interface{}) int {
	switch x := v.(type) {
	case string:
		return len(x)
	case []interface{}:
		return len(x)
	case []int:
		return len(x)
	case []float64:
		return len(x)
	case []string:
		return len(x)
	case []bool:
		return len(x)
	default:
		return 0
	}
}

func gts_typeof(v interface{}) string {
	if v == nil {
		return "null"
	}
	switch v.(type) {
	case float64:
		return "number"
	case string:
		return "string"
	case bool:
		return "boolean"
	default:
		return "object"
	}
}

func gts_tostring(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

func gts_toint(v interface{}) int {
	switch x := v.(type) {
	case int:
		return x
	case float64:
		return int(x)
	case string:
		var n int
		fmt.Sscanf(x, "%d", &n)
		return n
	case bool:
		if x {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func gts_tofloat(v interface{}) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case int:
		return float64(x)
	case string:
		var n float64
		fmt.Sscanf(x, "%f", &n)
		return n
	case bool:
		if x {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func gts_call(fn interface{}, args ...interface{}) interface{} {
	v := reflect.ValueOf(fn)
	in := make([]reflect.Value, len(args))
	fnType := v.Type()
	for i, arg := range args {
		if i < fnType.NumIn() {
			// Convert argument to expected type
			expectedType := fnType.In(i)
			argVal := reflect.ValueOf(arg)
			if argVal.Type().ConvertibleTo(expectedType) {
				in[i] = argVal.Convert(expectedType)
			} else {
				in[i] = argVal
			}
		} else {
			in[i] = reflect.ValueOf(arg)
		}
	}
	out := v.Call(in)
	if len(out) > 0 {
		return out[0].Interface()
	}
	return nil
}

func gts_tobool(v interface{}) bool {
	if v == nil {
		return false
	}
	switch x := v.(type) {
	case bool:
		return x
	case float64:
		return x != 0
	case string:
		return x != ""
	default:
		return true
	}
}

func gts_toarr_float(v []interface{}) []float64 {
	result := make([]float64, len(v))
	for i, x := range v {
		result[i] = gts_tofloat(x)
	}
	return result
}

func gts_toarr_int(v []interface{}) []int {
	result := make([]int, len(v))
	for i, x := range v {
		result[i] = gts_toint(x)
	}
	return result
}

func gts_map(arr interface{}, fn interface{}) interface{} {
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		return arr
	}
	resultType := reflect.SliceOf(v.Type().Elem())
	result := reflect.MakeSlice(resultType, 0, v.Len())
	f := reflect.ValueOf(fn)
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		out := f.Call([]reflect.Value{elem})
		result = reflect.Append(result, out[0])
	}
	return result.Interface()
}

func gts_filter(arr interface{}, fn interface{}) interface{} {
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		return arr
	}
	resultType := reflect.SliceOf(v.Type().Elem())
	result := reflect.MakeSlice(resultType, 0, v.Len())
	f := reflect.ValueOf(fn)
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		out := f.Call([]reflect.Value{elem})
		if out[0].Bool() {
			result = reflect.Append(result, elem)
		}
	}
	return result.Interface()
}

func gts_reduce(arr interface{}, initial interface{}, fn interface{}) interface{} {
	v := reflect.ValueOf(arr)
	f := reflect.ValueOf(fn)
	acc := reflect.ValueOf(initial)
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		out := f.Call([]reflect.Value{acc, elem})
		acc = out[0]
	}
	return acc.Interface()
}

func gts_find(arr interface{}, fn interface{}) interface{} {
	v := reflect.ValueOf(arr)
	f := reflect.ValueOf(fn)
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		out := f.Call([]reflect.Value{elem})
		if out[0].Bool() {
			return elem.Interface()
		}
	}
	var zero interface{}
	return zero
}

func gts_findIndex(arr interface{}, fn interface{}) int {
	v := reflect.ValueOf(arr)
	f := reflect.ValueOf(fn)
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		out := f.Call([]reflect.Value{elem})
		if out[0].Bool() {
			return i
		}
	}
	return -1
}

func gts_some(arr interface{}, fn interface{}) bool {
	v := reflect.ValueOf(arr)
	f := reflect.ValueOf(fn)
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		out := f.Call([]reflect.Value{elem})
		if out[0].Bool() {
			return true
		}
	}
	return false
}

func gts_every(arr interface{}, fn interface{}) bool {
	v := reflect.ValueOf(arr)
	f := reflect.ValueOf(fn)
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		out := f.Call([]reflect.Value{elem})
		if !out[0].Bool() {
			return false
		}
	}
	return true
}

// GTS_Promise represents a JavaScript-style Promise
type GTS_Promise[T any] struct {
	done    chan struct{}
	value   T
	err     error
	settled bool
}

func GTS_NewPromise[T any](executor func(resolve func(T), reject func(error))) *GTS_Promise[T] {
	p := &GTS_Promise[T]{done: make(chan struct{})}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					p.err = err
				} else {
					p.err = fmt.Errorf("%v", r)
				}
				if !p.settled {
					p.settled = true
					close(p.done)
				}
			}
		}()
		resolve := func(v T) {
			if !p.settled {
				p.value = v
				p.settled = true
				close(p.done)
			}
		}
		reject := func(e error) {
			if !p.settled {
				p.err = e
				p.settled = true
				close(p.done)
			}
		}
		executor(resolve, reject)
	}()
	return p
}

func gts_await[T any](p *GTS_Promise[T]) T {
	<-p.done
	if p.err != nil {
		panic(p.err)
	}
	return p.value
}

func GTS_Promise_Resolve[T any](value T) *GTS_Promise[T] {
	return GTS_NewPromise(func(resolve func(T), reject func(error)) {
		resolve(value)
	})
}

func GTS_Promise_Reject[T any](err error) *GTS_Promise[T] {
	return GTS_NewPromise(func(resolve func(T), reject func(error)) {
		reject(err)
	})
}

func GTS_Promise_All[T any](promises []*GTS_Promise[T]) *GTS_Promise[[]T] {
	return GTS_NewPromise(func(resolve func([]T), reject func(error)) {
		results := make([]T, len(promises))
		for i, p := range promises {
			<-p.done
			if p.err != nil {
				reject(p.err)
				return
			}
			results[i] = p.value
		}
		resolve(results)
	})
}

func GTS_Promise_Race[T any](promises []*GTS_Promise[T]) *GTS_Promise[T] {
	return GTS_NewPromise(func(resolve func(T), reject func(error)) {
		done := make(chan struct{})
		for _, p := range promises {
			go func(pr *GTS_Promise[T]) {
				<-pr.done
				select {
				case <-done:
					return
				default:
					close(done)
					if pr.err != nil {
						reject(pr.err)
					} else {
						resolve(pr.value)
					}
				}
			}(p)
		}
	})
}

type TreeNode struct {
	Value     int
	HasLeft   bool
	HasRight  bool
	LeftNode  *TreeNode
	RightNode *TreeNode
}

func NewTreeNode(value int) *TreeNode {
	this := &TreeNode{}
	this.Value = value
	this.HasLeft = false
	this.HasRight = false
	this.LeftNode = nil
	this.RightNode = nil
	return this
}

func (this *TreeNode) SetLeft(child *TreeNode) {
	this.LeftNode = child
	this.HasLeft = true
}

func (this *TreeNode) SetRight(child *TreeNode) {
	this.RightNode = child
	this.HasRight = true
}

func (this *TreeNode) Left() *TreeNode {
	var child *TreeNode = this.LeftNode
	if this.HasLeft && (child != nil) {
		return child
	}
	return this
}

func (this *TreeNode) Right() *TreeNode {
	var child *TreeNode = this.RightNode
	if this.HasRight && (child != nil) {
		return child
	}
	return this
}

func (this *TreeNode) IsLeaf() bool {
	return ((!this.HasLeft) && (!this.HasRight))
}

type BST struct {
	RootNode  *TreeNode
	NodeCount int
}

func NewBST() *BST {
	this := &BST{}
	this.RootNode = nil
	this.NodeCount = 0
	return this
}

func (this *BST) IsEmpty() bool {
	return (this.RootNode == nil)
}

func (this *BST) Root() *TreeNode {
	var r *TreeNode = this.RootNode
	if r != nil {
		return r
	}
	return NewTreeNode((-1))
}

func (this *BST) Insert(value int) {
	if this.RootNode == nil {
		this.RootNode = NewTreeNode(value)
		this.NodeCount = 1
	} else {
		this.InsertAt(this.Root(), value)
		this.NodeCount = (this.NodeCount + 1)
	}
}

func (this *BST) InsertAt(node *TreeNode, value int) {
	if value < node.Value {
		if node.HasLeft {
			this.InsertAt(node.Left(), value)
		} else {
			node.SetLeft(NewTreeNode(value))
		}
	} else if value > node.Value {
		if node.HasRight {
			this.InsertAt(node.Right(), value)
		} else {
			node.SetRight(NewTreeNode(value))
		}
	}
}

func (this *BST) Search(value int) bool {
	if this.IsEmpty() {
		return false
	}
	return this.SearchAt(this.Root(), value)
}

func (this *BST) SearchAt(node *TreeNode, value int) bool {
	if value == node.Value {
		return true
	}
	if (value < node.Value) && node.HasLeft {
		return this.SearchAt(node.Left(), value)
	}
	if (value > node.Value) && node.HasRight {
		return this.SearchAt(node.Right(), value)
	}
	return false
}

func (this *BST) FindMin() int {
	if this.IsEmpty() {
		return (-1)
	}
	return this.FindMinAt(this.Root())
}

func (this *BST) FindMinAt(node *TreeNode) int {
	if node.HasLeft {
		return this.FindMinAt(node.Left())
	}
	return node.Value
}

func (this *BST) FindMax() int {
	if this.IsEmpty() {
		return (-1)
	}
	return this.FindMaxAt(this.Root())
}

func (this *BST) FindMaxAt(node *TreeNode) int {
	if node.HasRight {
		return this.FindMaxAt(node.Right())
	}
	return node.Value
}

func (this *BST) Height() int {
	if this.IsEmpty() {
		return 0
	}
	return this.HeightAt(this.Root())
}

func (this *BST) HeightAt(node *TreeNode) int {
	var leftH int = 0
	var rightH int = 0
	if node.HasLeft {
		leftH = this.HeightAt(node.Left())
	}
	if node.HasRight {
		rightH = this.HeightAt(node.Right())
	}
	if leftH > rightH {
		return (1 + leftH)
	}
	return (1 + rightH)
}

func (this *BST) CountNodes() int {
	if this.IsEmpty() {
		return 0
	}
	return this.CountAt(this.Root())
}

func (this *BST) CountAt(node *TreeNode) int {
	var count int = 1
	if node.HasLeft {
		count = (count + this.CountAt(node.Left()))
	}
	if node.HasRight {
		count = (count + this.CountAt(node.Right()))
	}
	return count
}

func (this *BST) CountLeaves() int {
	if this.IsEmpty() {
		return 0
	}
	return this.CountLeavesAt(this.Root())
}

func (this *BST) CountLeavesAt(node *TreeNode) int {
	if node.IsLeaf() {
		return 1
	}
	var count int = 0
	if node.HasLeft {
		count = (count + this.CountLeavesAt(node.Left()))
	}
	if node.HasRight {
		count = (count + this.CountLeavesAt(node.Right()))
	}
	return count
}

func (this *BST) Sum() int {
	if this.IsEmpty() {
		return 0
	}
	return this.SumAt(this.Root())
}

func (this *BST) SumAt(node *TreeNode) int {
	var total int = node.Value
	if node.HasLeft {
		total = (total + this.SumAt(node.Left()))
	}
	if node.HasRight {
		total = (total + this.SumAt(node.Right()))
	}
	return total
}

func (this *BST) IsBalanced() bool {
	if this.IsEmpty() {
		return true
	}
	return (this.CheckBalanceAt(this.Root()) != (-1))
}

func (this *BST) CheckBalanceAt(node *TreeNode) int {
	var leftH int = 0
	var rightH int = 0
	if node.HasLeft {
		leftH = this.CheckBalanceAt(node.Left())
		if leftH == (-1) {
			return (-1)
		}
	}
	if node.HasRight {
		rightH = this.CheckBalanceAt(node.Right())
		if rightH == (-1) {
			return (-1)
		}
	}
	var diff int = (leftH - rightH)
	if diff < 0 {
		diff = (0 - diff)
	}
	if diff > 1 {
		return (-1)
	}
	if leftH > rightH {
		return (1 + leftH)
	}
	return (1 + rightH)
}

func (this *BST) IsValidBST() bool {
	if this.IsEmpty() {
		return true
	}
	return this.ValidateAt(this.Root(), (-2.147483647e+09), 2.147483647e+09)
}

func (this *BST) ValidateAt(node *TreeNode, minVal int, maxVal int) bool {
	if (node.Value <= minVal) || (node.Value >= maxVal) {
		return false
	}
	var leftValid bool = true
	var rightValid bool = true
	if node.HasLeft {
		leftValid = this.ValidateAt(node.Left(), minVal, node.Value)
	}
	if node.HasRight {
		rightValid = this.ValidateAt(node.Right(), node.Value, maxVal)
	}
	return (leftValid && rightValid)
}

func (this *BST) InorderTraversal() []int {
	var result []int = []int{}
	if !this.IsEmpty() {
		this.CollectInorder(this.Root(), result)
	}
	return result
}

func (this *BST) CollectInorder(node *TreeNode, result []int) {
	if node.HasLeft {
		this.CollectInorder(node.Left(), result)
	}
	func() int { result = append(result, node.Value); return len(result) }()
	if node.HasRight {
		this.CollectInorder(node.Right(), result)
	}
}

func (this *BST) PreorderTraversal() []int {
	var result []int = []int{}
	if !this.IsEmpty() {
		this.CollectPreorder(this.Root(), result)
	}
	return result
}

func (this *BST) CollectPreorder(node *TreeNode, result []int) {
	func() int { result = append(result, node.Value); return len(result) }()
	if node.HasLeft {
		this.CollectPreorder(node.Left(), result)
	}
	if node.HasRight {
		this.CollectPreorder(node.Right(), result)
	}
}

func (this *BST) PostorderTraversal() []int {
	var result []int = []int{}
	if !this.IsEmpty() {
		this.CollectPostorder(this.Root(), result)
	}
	return result
}

func (this *BST) CollectPostorder(node *TreeNode, result []int) {
	if node.HasLeft {
		this.CollectPostorder(node.Left(), result)
	}
	if node.HasRight {
		this.CollectPostorder(node.Right(), result)
	}
	func() int { result = append(result, node.Value); return len(result) }()
}

func (this *BST) Lca(a int, b int) int {
	if this.IsEmpty() {
		return (-1)
	}
	return this.FindLCA(this.Root(), a, b)
}

func (this *BST) FindLCA(node *TreeNode, a int, b int) int {
	if ((a < node.Value) && (b < node.Value)) && node.HasLeft {
		return this.FindLCA(node.Left(), a, b)
	}
	if ((a > node.Value) && (b > node.Value)) && node.HasRight {
		return this.FindLCA(node.Right(), a, b)
	}
	return node.Value
}

func (this *BST) HasPathSum(targetSum int) bool {
	if this.IsEmpty() {
		return false
	}
	return this.CheckPathSum(this.Root(), targetSum)
}

func (this *BST) CheckPathSum(node *TreeNode, remaining int) bool {
	var newRemaining int = (remaining - node.Value)
	if node.IsLeaf() {
		return (newRemaining == 0)
	}
	var leftHas bool = false
	var rightHas bool = false
	if node.HasLeft {
		leftHas = this.CheckPathSum(node.Left(), newRemaining)
	}
	if node.HasRight {
		rightHas = this.CheckPathSum(node.Right(), newRemaining)
	}
	return (leftHas || rightHas)
}

func (this *BST) GetLevel(level int) []int {
	var result []int = []int{}
	if !this.IsEmpty() {
		this.CollectLevel(this.Root(), level, 0, result)
	}
	return result
}

func (this *BST) CollectLevel(node *TreeNode, target int, current int, result []int) {
	if current == target {
		func() int { result = append(result, node.Value); return len(result) }()
	} else {
		if node.HasLeft {
			this.CollectLevel(node.Left(), target, (current + 1), result)
		}
		if node.HasRight {
			this.CollectLevel(node.Right(), target, (current + 1), result)
		}
	}
}

func (this *BST) LevelWidths() []int {
	var widths []int = []int{}
	var h int = this.Height()
	var i int = 0
	for i < h {
		var level []int = this.GetLevel(i)
		func() int { widths = append(widths, len(level)); return len(widths) }()
		i = (i + 1)
	}
	return widths
}

func (this *BST) Delete(value int) bool {
	if this.IsEmpty() {
		return false
	}
	if !this.Search(value) {
		return false
	}
	this.RootNode = this.DeleteAt(this.Root(), value)
	this.NodeCount = (this.NodeCount - 1)
	return true
}

func (this *BST) DeleteAt(node *TreeNode, value int) *TreeNode {
	if (value < node.Value) && node.HasLeft {
		var newLeft *TreeNode = this.DeleteAt(node.Left(), value)
		if newLeft != nil {
			node.SetLeft(newLeft)
		} else {
			node.HasLeft = false
			node.LeftNode = nil
		}
	} else if (value > node.Value) && node.HasRight {
		var newRight *TreeNode = this.DeleteAt(node.Right(), value)
		if newRight != nil {
			node.SetRight(newRight)
		} else {
			node.HasRight = false
			node.RightNode = nil
		}
	} else if value == node.Value {
		if (!node.HasLeft) && (!node.HasRight) {
			return nil
		}
		if !node.HasLeft {
			return node.RightNode
		}
		if !node.HasRight {
			return node.LeftNode
		}
		var successorVal int = this.FindMinAt(node.Right())
		node.Value = successorVal
		var newRight *TreeNode = this.DeleteAt(node.Right(), successorVal)
		if newRight != nil {
			node.SetRight(newRight)
		} else {
			node.HasRight = false
			node.RightNode = nil
		}
	}
	return node
}

func (this *BST) Diameter() int {
	if this.IsEmpty() {
		return 0
	}
	var maxDiam int = 0
	this.CalcDiameter(this.Root(), maxDiam)
	return maxDiam
}

func (this *BST) CalcDiameter(node *TreeNode, maxD int) int {
	var leftH int = 0
	var rightH int = 0
	if node.HasLeft {
		leftH = this.CalcDiameter(node.Left(), maxD)
	}
	if node.HasRight {
		rightH = this.CalcDiameter(node.Right(), maxD)
	}
	var currDiam int = (leftH + rightH)
	if currDiam > maxD {
		maxD = currDiam
	}
	if leftH > rightH {
		return (1 + leftH)
	}
	return (1 + rightH)
}

type RNG struct {
	State int
}

func NewRNG(seed int) *RNG {
	this := &RNG{}
	this.State = seed
	return this
}

func (this *RNG) Next() int {
	var a int = 1.103515245e+09
	var c int = 12345
	var m int = 2.147483647e+09
	this.State = (((a * this.State) + c) % m)
	if this.State < 0 {
		this.State = (0 - this.State)
	}
	return this.State
}

func (this *RNG) NextInRange(maxVal int) int {
	var v int = (this.Next() % maxVal)
	if v < 0 {
		v = (0 - v)
	}
	return v
}

func boolStr(val bool) string {
	if val {
		return "true"
	}
	return "false"
}

func benchmark(n int) {
	fmt.Println((("=== Binary Tree Benchmark (" + gts_tostring(n)) + " elements) ==="))
	fmt.Println("")
	var tree *BST = NewBST()
	var rng *RNG = NewRNG(42)
	fmt.Println("1. Inserting elements...")
	var i int = 0
	for i < n {
		var v int = rng.NextInRange(100000)
		tree.Insert(v)
		i = (i + 1)
	}
	fmt.Println(("   Tree size: " + gts_tostring(tree.CountNodes())))
	fmt.Println("")
	fmt.Println("2. Tree Properties:")
	fmt.Println(("   Height: " + gts_tostring(tree.Height())))
	fmt.Println(("   Leaves: " + gts_tostring(tree.CountLeaves())))
	fmt.Println(("   Balanced: " + boolStr(tree.IsBalanced())))
	fmt.Println(("   Valid BST: " + boolStr(tree.IsValidBST())))
	fmt.Println(("   Sum: " + gts_tostring(tree.Sum())))
	fmt.Println(("   Min: " + gts_tostring(tree.FindMin())))
	fmt.Println(("   Max: " + gts_tostring(tree.FindMax())))
	fmt.Println("")
	fmt.Println("3. Level widths (first 10):")
	var widths []int = tree.LevelWidths()
	var j int = 0
	for (j < 10) && (j < len(widths)) {
		fmt.Println(((("   Level " + gts_tostring(j)) + ": ") + gts_tostring(widths[int(j)])))
		j = (j + 1)
	}
	fmt.Println("")
	fmt.Println("4. Searching 100 random values...")
	var found int = 0
	var searchRng *RNG = NewRNG(12345)
	var k int = 0
	for k < 100 {
		var sv int = searchRng.NextInRange(100000)
		if tree.Search(sv) {
			found = (found + 1)
		}
		k = (k + 1)
	}
	fmt.Println((("   Found: " + gts_tostring(found)) + "/100"))
	fmt.Println("")
	fmt.Println("5. LCA of min and max:")
	var minV int = tree.FindMin()
	var maxV int = tree.FindMax()
	var lcaV int = tree.Lca(minV, maxV)
	fmt.Println(((((("   LCA(" + gts_tostring(minV)) + ",") + gts_tostring(maxV)) + ") = ") + gts_tostring(lcaV)))
	fmt.Println("")
	fmt.Println("6. Path sum check:")
	var nc int = tree.CountNodes()
	var avgH int = tree.Height()
	var testS int = (gts_toint((tree.Sum() / nc)) * avgH)
	fmt.Println(((("   Has path sum " + gts_tostring(testS)) + ": ") + boolStr(tree.HasPathSum(testS))))
	fmt.Println("")
	fmt.Println("7. Traversals (first 5 elements):")
	var inord []int = tree.InorderTraversal()
	var preord []int = tree.PreorderTraversal()
	var postord []int = tree.PostorderTraversal()
	fmt.Println("   Inorder:   ")
	var m int = 0
	for (m < 5) && (m < len(inord)) {
		fmt.Println((gts_tostring(inord[int(m)]) + " "))
		m = (m + 1)
	}
	fmt.Println("")
	fmt.Println("   Preorder:  ")
	m = 0
	for (m < 5) && (m < len(preord)) {
		fmt.Println((gts_tostring(preord[int(m)]) + " "))
		m = (m + 1)
	}
	fmt.Println("")
	fmt.Println("   Postorder: ")
	var plen int = len(postord)
	var si int = (plen - 5)
	if si < 0 {
		si = 0
	}
	m = si
	for m < plen {
		fmt.Println((gts_tostring(postord[int(m)]) + " "))
		m = (m + 1)
	}
	fmt.Println("")
	fmt.Println("")
	fmt.Println("8. Verify BST property:")
	var sorted bool = true
	var s int = 1
	for s < len(inord) {
		if inord[int(s)] < inord[int((s-1))] {
			sorted = false
		}
		s = (s + 1)
	}
	fmt.Println(("   Inorder sorted: " + boolStr(sorted)))
	fmt.Println("")
	fmt.Println("9. Deleting 50 random elements...")
	var deleted int = 0
	var delRng *RNG = NewRNG(99999)
	var d int = 0
	for d < 50 {
		var dv int = delRng.NextInRange(100000)
		if tree.Delete(dv) {
			deleted = (deleted + 1)
		}
		d = (d + 1)
	}
	fmt.Println(("   Deleted: " + gts_tostring(deleted)))
	fmt.Println(("   New size: " + gts_tostring(tree.CountNodes())))
	fmt.Println(("   Still valid: " + boolStr(tree.IsValidBST())))
	fmt.Println("")
	fmt.Println("10. Post-deletion metrics:")
	fmt.Println(("   Height: " + gts_tostring(tree.Height())))
	fmt.Println(("   Leaves: " + gts_tostring(tree.CountLeaves())))
	var level5 []int = tree.GetLevel(5)
	fmt.Println(("   Nodes at level 5: " + gts_tostring(len(level5))))
	fmt.Println("")
	fmt.Println("=== Benchmark Complete ===")
}

func main() {
	fmt.Println("Binary Tree Operations Benchmark")
	fmt.Println("================================")
	fmt.Println("")
	benchmark(1e+06)
}
