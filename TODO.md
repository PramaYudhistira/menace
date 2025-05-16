* Restructure this spaghetti...
    * Things like 2 definitions for Model struct shouldn't be there

* Processing message should be more elegant animations wise
    * The processing message only shows up when the llm response is received, therefore meaningless for now.
    * Check if its anything to do with mutex locks
