#!/bin/bash

# Go Laboratory Skill Assessment Tool
# Evaluates learner progress and provides personalized recommendations

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Assessment state
declare -A SCORES
declare -A COMPLETED_MODULES
TOTAL_SCORE=0
MAX_SCORE=0

# Print colored output
print_header() {
    echo -e "${PURPLE}=====================================${NC}"
    echo -e "${PURPLE}  Go Laboratory Skill Assessment${NC}"
    echo -e "${PURPLE}=====================================${NC}"
    echo
}

print_section() {
    echo -e "${CYAN}$1${NC}"
    echo -e "${CYAN}$(echo "$1" | sed 's/./-/g')${NC}"
}

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[✓]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[!]${NC} $1"
}

print_error() {
    echo -e "${RED}[✗]${NC} $1"
}

print_question() {
    echo -e "${CYAN}$1${NC}"
}

# Check if directory exists and has expected content
check_directory() {
    local dir=$1
    local description=$2
    
    if [[ -d "$dir" ]]; then
        print_success "$description found"
        return 0
    else
        print_error "$description missing"
        return 1
    fi
}

# Check if Go files in directory compile
check_compilation() {
    local dir=$1
    local description=$2
    
    if [[ ! -d "$dir" ]]; then
        print_error "$description directory not found"
        return 1
    fi
    
    local go_files=$(find "$dir" -name "*.go" -type f | wc -l)
    if [[ $go_files -eq 0 ]]; then
        print_warning "$description has no Go files"
        return 0
    fi
    
    cd "$dir" 2>/dev/null || return 1
    
    if go build ./... >/dev/null 2>&1; then
        print_success "$description code compiles"
        return 0
    else
        print_error "$description code has compilation errors"
        return 1
    fi
}

# Check if tests exist and pass
check_tests() {
    local dir=$1
    local description=$2
    
    if [[ ! -d "$dir" ]]; then
        return 1
    fi
    
    cd "$dir" 2>/dev/null || return 1
    
    local test_files=$(find . -name "*_test.go" -type f | wc -l)
    if [[ $test_files -eq 0 ]]; then
        print_warning "$description has no tests"
        return 0
    fi
    
    if go test ./... >/dev/null 2>&1; then
        print_success "$description tests pass"
        return 0
    else
        print_error "$description tests are failing"
        return 1
    fi
}

# Assess fundamentals completion
assess_fundamentals() {
    print_section "📚 Fundamentals Assessment"
    
    local score=0
    local max_score=30
    
    # Check hello world
    if check_directory "learning/01-fundamentals/hello" "Hello World example"; then
        ((score += 5))
        if check_compilation "learning/01-fundamentals/hello" "Hello World"; then
            ((score += 5))
        fi
    fi
    
    # Check math module
    if check_directory "learning/01-fundamentals/math" "Math examples"; then
        ((score += 5))
        if check_compilation "learning/01-fundamentals/math" "Math examples"; then
            ((score += 5))
        fi
        if check_tests "learning/01-fundamentals/math" "Math examples"; then
            ((score += 10))
        fi
    fi
    
    SCORES["fundamentals"]=$score
    ((MAX_SCORE += max_score))
    ((TOTAL_SCORE += score))
    
    echo "Fundamentals Score: $score/$max_score"
    
    if [[ $score -ge 25 ]]; then
        print_success "Fundamentals mastery achieved!"
        COMPLETED_MODULES["fundamentals"]=true
    elif [[ $score -ge 15 ]]; then
        print_warning "Fundamentals partially completed"
    else
        print_error "Fundamentals need attention"
    fi
    
    echo
}

# Assess intermediate completion
assess_intermediate() {
    print_section "🎯 Intermediate Assessment"
    
    local score=0
    local max_score=40
    
    # Check composition examples
    if check_directory "learning/02-intermediate/composition" "Composition examples"; then
        ((score += 10))
        if check_compilation "learning/02-intermediate/composition" "Composition"; then
            ((score += 10))
        fi
    fi
    
    # Check HTTP services
    if check_directory "examples/02-intermediate/http-services" "HTTP services"; then
        ((score += 5))
        if check_compilation "examples/02-intermediate/http-services/jsonplaceholder" "JSONPlaceholder client"; then
            ((score += 5))
        fi
        if check_compilation "examples/02-intermediate/http-services/server" "HTTP server"; then
            ((score += 5))
        fi
        if check_tests "examples/02-intermediate/http-services" "HTTP services"; then
            ((score += 5))
        fi
    fi
    
    SCORES["intermediate"]=$score
    ((MAX_SCORE += max_score))
    ((TOTAL_SCORE += score))
    
    echo "Intermediate Score: $score/$max_score"
    
    if [[ $score -ge 35 ]]; then
        print_success "Intermediate mastery achieved!"
        COMPLETED_MODULES["intermediate"]=true
    elif [[ $score -ge 20 ]]; then
        print_warning "Intermediate partially completed"
    else
        print_error "Intermediate needs work"
    fi
    
    echo
}

# Assess advanced completion
assess_advanced() {
    print_section "🚀 Advanced Assessment"
    
    local score=0
    local max_score=50
    
    # Check concurrency
    if check_directory "learning/03-advanced/concurrency" "Concurrency examples"; then
        ((score += 10))
        if check_compilation "learning/03-advanced/concurrency" "Concurrency"; then
            ((score += 5))
        fi
    fi
    
    # Check patterns
    if check_directory "learning/03-advanced/patterns" "Design patterns"; then
        ((score += 10))
        if check_compilation "learning/03-advanced/patterns/creational/builder" "Builder pattern"; then
            ((score += 5))
        fi
        if check_compilation "learning/03-advanced/patterns/creational/singleton" "Singleton pattern"; then
            ((score += 5))
        fi
    fi
    
    # Check production API
    if check_directory "examples/03-advanced/production-api" "Production API"; then
        ((score += 10))
        if check_compilation "examples/03-advanced/production-api" "Production API"; then
            ((score += 5))
        fi
    fi
    
    SCORES["advanced"]=$score
    ((MAX_SCORE += max_score))
    ((TOTAL_SCORE += score))
    
    echo "Advanced Score: $score/$max_score"
    
    if [[ $score -ge 40 ]]; then
        print_success "Advanced mastery achieved!"
        COMPLETED_MODULES["advanced"]=true
    elif [[ $score -ge 25 ]]; then
        print_warning "Advanced partially completed"
    else
        print_error "Advanced needs significant work"
    fi
    
    echo
}

# Assess testing knowledge
assess_testing() {
    print_section "🧪 Testing Assessment"
    
    local score=0
    local max_score=30
    
    # Check calculator examples
    if check_directory "examples/01-beginner/calculator" "Calculator examples"; then
        ((score += 5))
        if check_tests "examples/01-beginner/calculator/v1" "Calculator v1"; then
            ((score += 5))
        fi
        if check_tests "examples/01-beginner/calculator/v2" "Calculator v2"; then
            ((score += 10))
        fi
    fi
    
    # Check testing directory
    if check_directory "testing" "Testing examples"; then
        ((score += 5))
        if check_tests "testing/mocking" "Mocking examples"; then
            ((score += 5))
        fi
    fi
    
    SCORES["testing"]=$score
    ((MAX_SCORE += max_score))
    ((TOTAL_SCORE += score))
    
    echo "Testing Score: $score/$max_score"
    
    if [[ $score -ge 25 ]]; then
        print_success "Testing mastery achieved!"
        COMPLETED_MODULES["testing"]=true
    elif [[ $score -ge 15 ]]; then
        print_warning "Testing partially completed"
    else
        print_error "Testing needs attention"
    fi
    
    echo
}

# Generate skill level assessment
assess_skill_level() {
    print_section "📊 Overall Skill Assessment"
    
    local percentage=$((TOTAL_SCORE * 100 / MAX_SCORE))
    
    echo "Overall Score: $TOTAL_SCORE/$MAX_SCORE ($percentage%)"
    echo
    
    if [[ $percentage -ge 90 ]]; then
        print_success "🏆 Expert Level (90%+)"
        echo "You have mastered Go development!"
        echo "Consider:"
        echo "  • Contributing to open source projects"
        echo "  • Mentoring other developers"
        echo "  • Leading technical architecture decisions"
        echo "  • Exploring cutting-edge Go features"
    elif [[ $percentage -ge 70 ]]; then
        print_success "🎯 Advanced Level (70-89%)"
        echo "You have strong Go skills!"
        echo "Focus on:"
        echo "  • Production deployment patterns"
        echo "  • Performance optimization"
        echo "  • Advanced architectural patterns"
        echo "  • Building complex distributed systems"
    elif [[ $percentage -ge 50 ]]; then
        print_warning "📈 Intermediate Level (50-69%)"
        echo "You're making good progress!"
        echo "Next steps:"
        echo "  • Complete remaining intermediate topics"
        echo "  • Build more complex HTTP services"
        echo "  • Master testing strategies"
        echo "  • Learn production best practices"
    elif [[ $percentage -ge 25 ]]; then
        print_warning "📚 Beginner+ Level (25-49%)"
        echo "You have the basics down!"
        echo "Continue with:"
        echo "  • Strengthening Go fundamentals"
        echo "  • Building simple projects"
        echo "  • Learning about interfaces and composition"
        echo "  • Practicing with examples"
    else
        print_error "🌱 Beginner Level (0-24%)"
        echo "Great start on your Go journey!"
        echo "Focus on:"
        echo "  • Go syntax and basic concepts"
        echo "  • Running and modifying examples"
        echo "  • Understanding functions and error handling"
        echo "  • Following the learning path step by step"
    fi
    
    echo
}

# Provide personalized recommendations
provide_recommendations() {
    print_section "💡 Personalized Recommendations"
    
    # Recommendations based on completed modules
    if [[ "${COMPLETED_MODULES[fundamentals]}" != "true" ]]; then
        echo "🔸 Complete fundamentals first:"
        echo "   • Work through learning/01-fundamentals/"
        echo "   • Ensure all examples compile and run"
        echo "   • Understand error handling patterns"
        echo
    fi
    
    if [[ "${COMPLETED_MODULES[fundamentals]}" == "true" && "${COMPLETED_MODULES[intermediate]}" != "true" ]]; then
        echo "🔸 Focus on intermediate concepts:"
        echo "   • Study interfaces and composition"
        echo "   • Build HTTP services from examples"
        echo "   • Practice with JSON handling"
        echo "   • Learn service architecture patterns"
        echo
    fi
    
    if [[ "${COMPLETED_MODULES[intermediate]}" == "true" && "${COMPLETED_MODULES[advanced]}" != "true" ]]; then
        echo "🔸 Advance to production patterns:"
        echo "   • Master concurrency with goroutines and channels"
        echo "   • Implement design patterns"
        echo "   • Study the production API example"
        echo "   • Learn deployment and monitoring"
        echo
    fi
    
    if [[ "${COMPLETED_MODULES[testing]}" != "true" ]]; then
        echo "🔸 Strengthen testing skills:"
        echo "   • Practice table-driven tests"
        echo "   • Learn BDD with Ginkgo/Gomega"
        echo "   • Master mocking patterns"
        echo "   • Write integration tests"
        echo
    fi
    
    # General recommendations
    echo "📋 General next steps:"
    echo "  1. Follow the LEARNING_ROADMAP.md for structured progression"
    echo "  2. Use INDEX.md to find specific concepts quickly"
    echo "  3. Try building projects with templates/"
    echo "  4. Join the Go community and contribute to projects"
    echo
}

# Check environment and tools
check_environment() {
    print_section "🔧 Environment Check"
    
    # Check Go installation
    if command -v go &> /dev/null; then
        local go_version=$(go version | awk '{print $3}')
        print_success "Go installed: $go_version"
    else
        print_error "Go not installed or not in PATH"
        echo "Install Go from: https://golang.org/dl/"
        return 1
    fi
    
    # Check development tools
    local tools=("golangci-lint" "gofmt" "goimports")
    for tool in "${tools[@]}"; do
        if command -v "$tool" &> /dev/null; then
            print_success "$tool available"
        else
            print_warning "$tool not found (recommended for development)"
        fi
    done
    
    echo
    return 0
}

# Interactive mode for deeper assessment
interactive_assessment() {
    print_section "🤔 Interactive Skills Check"
    
    echo "Answer these questions to get more personalized recommendations:"
    echo
    
    # Experience level
    print_question "1. What's your overall programming experience?"
    echo "   a) New to programming"
    echo "   b) Some programming experience in other languages" 
    echo "   c) Experienced programmer, new to Go"
    echo "   d) Experienced with Go"
    read -p "   Your answer (a/b/c/d): " experience
    
    # Learning goals
    print_question "2. What's your primary learning goal?"
    echo "   a) Learn programming fundamentals"
    echo "   b) Build web APIs and services"
    echo "   c) System programming and CLI tools"
    echo "   d) Contribute to Go projects"
    read -p "   Your answer (a/b/c/d): " goal
    
    # Time availability
    print_question "3. How much time can you dedicate to learning?"
    echo "   a) 1-2 hours per week"
    echo "   b) 3-5 hours per week"
    echo "   c) 10+ hours per week"
    echo "   d) Full-time learning"
    read -p "   Your answer (a/b/c/d): " time
    
    echo
    print_section "📝 Customized Learning Plan"
    
    # Provide customized advice based on answers
    case $experience in
        a)
            echo "🌱 New Programmer Path:"
            echo "  • Start with fundamentals and take your time"
            echo "  • Focus on understanding concepts deeply"
            echo "  • Don't rush to advanced topics"
            ;;
        b)
            echo "🔄 Language Transfer Path:"
            echo "  • Compare Go concepts to languages you know"
            echo "  • Focus on Go-specific features (goroutines, interfaces)"
            echo "  • Practice idiomatic Go patterns"
            ;;
        c)
            echo "⚡ Go-Specific Path:"
            echo "  • Focus on Go's unique features and philosophy"
            echo "  • Study production patterns and best practices"
            echo "  • Compare with languages you know"
            ;;
        d)
            echo "🏆 Mastery Path:"
            echo "  • Focus on advanced patterns and performance"
            echo "  • Contribute to this repository"
            echo "  • Mentor others in the community"
            ;;
    esac
    
    case $goal in
        a)
            echo
            echo "🎯 Focus Areas for Programming Fundamentals:"
            echo "  • Complete all fundamentals thoroughly"
            echo "  • Build simple CLI tools"
            echo "  • Master testing from the beginning"
            ;;
        b)
            echo
            echo "🌐 Focus Areas for Web Development:"
            echo "  • Master HTTP services examples"
            echo "  • Study the production API architecture"
            echo "  • Learn deployment and monitoring patterns"
            ;;
        c)
            echo
            echo "⚙️ Focus Areas for System Programming:"
            echo "  • Master concurrency patterns"
            echo "  • Build CLI tools with Cobra"
            echo "  • Study performance optimization"
            ;;
        d)
            echo
            echo "🤝 Focus Areas for Open Source:"
            echo "  • Master all repository content"
            echo "  • Study contribution guidelines"
            echo "  • Start with documentation improvements"
            ;;
    esac
    
    case $time in
        a)
            echo
            echo "⏰ Recommended Schedule (1-2 hours/week):"
            echo "  • 6+ months for fundamentals"
            echo "  • 12+ months for intermediate"
            echo "  • Focus on depth over speed"
            ;;
        b)
            echo
            echo "⏰ Recommended Schedule (3-5 hours/week):"
            echo "  • 2-3 months for fundamentals"
            echo "  • 4-6 months for intermediate"
            echo "  • 8-12 months for advanced mastery"
            ;;
        c|d)
            echo
            echo "⏰ Recommended Schedule (10+ hours/week):"
            echo "  • 2-4 weeks for fundamentals"
            echo "  • 6-8 weeks for intermediate"
            echo "  • 3-4 months for advanced mastery"
            ;;
    esac
    
    echo
}

# Main execution
main() {
    print_header
    
    print_info "Assessing your Go Laboratory progress..."
    echo
    
    # Check if we're in the right directory
    if [[ ! -f "INDEX.md" ]] || [[ ! -d "learning" ]]; then
        print_error "Please run this script from the go-lab repository root"
        exit 1
    fi
    
    # Environment check
    if ! check_environment; then
        print_error "Environment issues detected. Please fix before continuing."
        exit 1
    fi
    
    # Run assessments
    assess_fundamentals
    assess_intermediate
    assess_advanced
    assess_testing
    
    # Overall assessment
    assess_skill_level
    
    # Recommendations
    provide_recommendations
    
    # Interactive mode
    echo "Would you like to take an interactive assessment for personalized recommendations? (y/n)"
    read -p "Answer: " interactive
    
    if [[ "$interactive" =~ ^[Yy] ]]; then
        echo
        interactive_assessment
    fi
    
    # Final message
    print_section "🎉 Assessment Complete"
    echo "Thank you for using the Go Laboratory Skill Assessment!"
    echo
    echo "Next steps:"
    echo "• Review your scores and recommendations above"
    echo "• Follow the suggested learning path"
    echo "• Re-run this assessment periodically to track progress"
    echo "• Share your progress with the community"
    echo
    print_success "Happy learning! 🚀"
}

# Execute main function
main "$@"