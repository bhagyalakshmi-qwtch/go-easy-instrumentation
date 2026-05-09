instrument.go
1. instrumentCmd — wired up the --exclude flag


Before:

Run: func(cmd *cobra.Command, args []string) {
    Instrument(args[0])
},


After:

Run: func(cmd *cobra.Command, args []string) {
    exclusions := parseExcludeDirs(excludeDirs)
    Instrument(args[0], exclusions...)
},


The flag value was being ignored before. Now it's parsed and passed to Instrument.

##############

2. parseExcludeDirs — new helper function

func parseExcludeDirs(raw string) []string {
    if raw == "" {
        return nil
    }
    var result []string
    for _, dir := range strings.Split(raw, ",") {
        if trimmed := strings.TrimSpace(dir); trimmed != "" {
            result = append(result, trimmed)
        }
    }
    return result
}


Splits "helper,models,repository" → ["helper", "models", "repository"]. Previously this logic was duplicated inside interactive.go. Now it lives in one place and is reused everywhere.

#################

3. buildLoadPatterns — new core function (the main fix)


func buildLoadPatterns(packagePath string, exclusions []string) ([]string, error) {
    if len(exclusions) == 0 {
        return []string{"./..."}, nil  // no change in behavior
    }
    // reads top-level directories, skips excluded ones
    // always includes "." for main.go
    patterns := []string{"."}
    for _, entry := range entries {
        if !entry.IsDir() { continue }
        if excluded[entry.Name()] { continue }
        patterns = append(patterns, "./"+entry.Name()+"/...")
    }
    return patterns, nil
}


This is the heart of the fix. Instead of always loading ./... (everything), it builds explicit patterns like [".", "./controller/...", "./service/..."] — skipping excluded directories.

#####################


4. Instrument — signature changed


Before:

func Instrument(packagePath string, patterns ...string) {


After:

func Instrument(packagePath string, exclusions ...string) {
The parameter was renamed from patterns to exclusions to reflect what's actually being passed now.

########################

5. runTextMode — now resolves exclusions to patterns


Before:

func runTextMode(packagePath string, patterns []string, outputFile string) {
    // passed patterns directly to instrumentPackages
}


After:


func runTextMode(packagePath string, exclusions []string, outputFile string) {
    patterns, err := buildLoadPatterns(packagePath, exclusions) // resolve here
    instrumentPackages(packagePath, patterns, outputFile)
}

################

6. runTUIMode — same as runTextMode
Before:

func runTUIMode(packagePath string, patterns []string, outputFile string) {
    loadPatterns := patterns
    if len(loadPatterns) == 0 {
        loadPatterns = []string{"./..."}
    }
    decorator.Load(..., loadPatterns...)
}
After:

func runTUIMode(packagePath string, exclusions []string, outputFile string) {
    patterns, err := buildLoadPatterns(packagePath, exclusions) // resolve here
    decorator.Load(..., patterns...)
}

#####################

interactive.go
7. Removed duplicated parsing logic, fixed what gets passed to Instrument
Before:

// duplicated split/trim logic
var exclusions []string
if excludeDirs != "" {
    for _, dir := range strings.Split(excludeDirs, ",") {
        ...
    }
}
// was passing file paths — wrong!
Instrument(".", files...)
After:

// reuses the shared helper
exclusions := parseExcludeDirs(excludeDirs)
// now correctly passes exclusions
Instrument(".", exclusions...)
The old code was passing individual .go file paths as patterns to Instrument, which was incorrect. Now it passes exclusion directory names, consistent with how instrument subcommand works.