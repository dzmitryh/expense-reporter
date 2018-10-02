Expense reporter
==========================

The idea behind is to create target version of report with
stops that specified at the beginning of main package.

#### Prerequisites:
* In order to be able to build this project you need to have GO :  
> *Download and install go from [golang.org](https://golang.org/dl/)*
* Set environment variable with your bus stops: `export REPORTING_BUS_STOPS="Schiphol Airport | Schiphol-Rijk, Boeingavenue | Schiphol-Rijk, Beechavenue"`
> Pay attention that delimiter is "|" symbol
#### How to build/run:
* go build
* expense-reporter {inputFilePath/FileName.csv} {outputDirPath}