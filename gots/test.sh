#!/bin/bash

# goTS Test Runner
# Usage: ./test.sh [unit|integration|all]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Build the binary first
echo -e "${YELLOW}Building gots binary...${NC}"
go build -o gots ./cmd/gots
echo -e "${GREEN}✓ Build successful${NC}\n"

run_unit_tests() {
    echo -e "${YELLOW}Running unit tests...${NC}"
    go test -v ./...
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Unit tests passed${NC}\n"
        return 0
    else
        echo -e "${RED}✗ Unit tests failed${NC}\n"
        return 1
    fi
}

run_integration_tests() {
    echo -e "${YELLOW}Running integration tests...${NC}"

    # Find all .gts files in test directory
    test_files=(test/*.gts)
    passed=0
    failed=0

    for file in "${test_files[@]}"; do
        if [ -f "$file" ]; then
            echo -e "\n${YELLOW}Testing: $file${NC}"
            if ./gots run "$file"; then
                echo -e "${GREEN}✓ $file passed${NC}"
                ((passed++))
            else
                echo -e "${RED}✗ $file failed${NC}"
                ((failed++))
            fi
        fi
    done

    echo -e "\n${YELLOW}Integration test results:${NC}"
    echo -e "  Passed: ${GREEN}$passed${NC}"
    echo -e "  Failed: ${RED}$failed${NC}"

    if [ $failed -eq 0 ]; then
        echo -e "${GREEN}✓ All integration tests passed${NC}\n"
        return 0
    else
        echo -e "${RED}✗ Some integration tests failed${NC}\n"
        return 1
    fi
}

# Parse command line argument
case "${1:-all}" in
    unit)
        run_unit_tests
        ;;
    integration)
        run_integration_tests
        ;;
    all)
        unit_result=0
        integration_result=0

        run_unit_tests || unit_result=$?
        run_integration_tests || integration_result=$?

        if [ $unit_result -eq 0 ] && [ $integration_result -eq 0 ]; then
            echo -e "${GREEN}✓ All tests passed!${NC}"
            exit 0
        else
            echo -e "${RED}✗ Some tests failed${NC}"
            exit 1
        fi
        ;;
    *)
        echo "Usage: $0 [unit|integration|all]"
        echo ""
        echo "  unit         - Run Go unit tests only"
        echo "  integration  - Run .gts integration tests only"
        echo "  all          - Run both unit and integration tests (default)"
        exit 1
        ;;
esac
