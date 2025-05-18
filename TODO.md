* Restructure this spaghetti...
    * Things like 2 definitions for Model struct shouldn't be there
    * Note: as opposed to the deprecated handling method, we can check for new lines and therefore travel up and down new lines based on terminal width. Perhaps with a new field to handle instead of raw lines have "displayedLines" field(needs further investigation)

* Processing message should be more elegant animations wise
    * The processing message only shows up when the llm response is received, therefore meaningless for now.
    * Check if its anything to do with mutex locks
