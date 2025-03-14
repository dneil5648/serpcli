package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func GenerateDork(kw, dork string) string {
	return fmt.Sprintf(`"%v" AND %v`, kw, dork)

}

var allDorks []string = []string{
	// File Types
	"filetype:xls", "filetype:xlsx", "filetype:doc", "filetype:docx", "filetype:ppt", "filetype:pptx", "filetype:pdf", "filetype:csv", "filetype:txt", "filetype:rtf", "filetype:odt", "filetype:ods", "filetype:odp", "filetype:xml", "filetype:json", "filetype:yaml", "filetype:yml", "filetype:ini", "filetype:cfg", "filetype:conf", "filetype:log", "filetype:sql", "filetype:db", "filetype:dbf", "filetype:mdb", "filetype:accdb", "filetype:sqlite", "filetype:tar", "filetype:gz", "filetype:zip", "filetype:rar", "filetype:7z", "filetype:bak", "filetype:backup", "filetype:bkf", "filetype:bkp", "filetype:iso", "filetype:img", "filetype:vmdk", "filetype:vdi", "filetype:ova", "filetype:ovf", "filetype:pem", "filetype:key", "filetype:crt", "filetype:cert", "filetype:p12", "filetype:pfx", "filetype:der",

	// Code Types
	"filetype:py", "filetype:pyc", "filetype:pyd", "filetype:pyo", "filetype:pyw", "filetype:pyz", "filetype:php", "filetype:phps", "filetype:php3", "filetype:php4", "filetype:php5", "filetype:php7", "filetype:phtml", "filetype:js", "filetype:jsx", "filetype:ts", "filetype:tsx", "filetype:coffee", "filetype:litcoffee", "filetype:dart", "filetype:go", "filetype:gohtml", "filetype:sh", "filetype:bash", "filetype:zsh", "filetype:pl", "filetype:pm", "filetype:psm1", "filetype:ps1", "filetype:ps1xml", "filetype:psc1", "filetype:pssc", "filetype:c", "filetype:cpp", "filetype:cs", "filetype:csx", "filetype:h", "filetype:hpp", "filetype:hxx", "filetype:java", "filetype:class", "filetype:jar", "filetype:jsp", "filetype:aspx", "filetype:asp", "filetype:asm", "filetype:s", "filetype:swift", "filetype:sqlite", "filetype:sql", "filetype:pgsql", "filetype:plsql", "filetype:mongodb", "filetype:perl", "filetype:rb", "filetype:erb", "filetype:html", "filetype:htm", "filetype:css", "filetype:scss", "filetype:sass", "filetype:less", "filetype:vue", "filetype:xml", "filetype:yml", "filetype:yaml", "filetype:json", "filetype:csv", "filetype:log", "filetype:txt", "filetype:md", "filetype:markdown", "filetype:rst", "filetype:tex", "filetype:bib", "filetype:ods", "filetype:xls", "filetype:xlsx", "filetype:doc", "filetype:docx", "filetype:ppt", "filetype:pptx", "filetype:key", "filetype:pub", "filetype:crt", "filetype:pem", "filetype:asc", "filetype:ppk", "filetype:cer", "filetype:pfx", "filetype:p12",

	// File Servers
	"intitle:\"index of \" \"parent directory\"", "intitle:\"index of\" inurl:ftp", "intitle:\"index of\" inurl:webdav",

	// Databases
	"filetype:sql \"dump\"", "filetype:cnf", "filetype:conf", "filetype:cfg mysql", "filetype:json \"mongodb\"", "intitle:\"MongoDB\" \"database\"", "intitle:\"index of\" /_utils/ \"CouchDB\"", "intitle:\"index of\" /_search", "filetype:conf \"postgresql\"", "filetype:sql \"mysql dump\"", "mariadb dump", "intitle:\"phpMyAdmin\" \"Welcome to phpMyAdmin\"", "filetype:bak", "filetype:backup", "filetype:sql",

	// Code Sites
	"site:github.com", "site:raw.githubusercontent.com", "site:gitlab.com", "site:bitbucket.org", "site:sourceforge.net", "site:codepen.io", "site:jsfiddle.net", "site:pastebin.com", "site:repl.it", "site:gist.github.com", "site:launchpad.net", "site:code.google.com", "site:codeplex.com", "site:jsdelivr.com", "site:npmjs.com", "site:pypi.org", "site:rubygems.org", "site:packagist.org", "site:maven.org", "site:nuget.org", "site:apache.org/dist", "site:cran.r-project.org/src/contrib", "site:cpan.org", "site:ctan.org", "site:perforce.com",

	// Cloud Storage Sites
	"site:amazonaws.com inurl:s3", "site:digitaloceanspaces.com", "site:wasabisys.com", "site:backblazeb2.com", "site:dream.io", "site:rackspacecloud.com", "site:scw.cloud", "site:vultrobjects.com", "site:oraclecloud.com inurl:objectstorage", "site:cloud.ibm.com in:objectstorage", "site:storage.googleapis.com", "site:blob.core.windows.net", "site:aliyuncs.com", "site:alibabacloud.com",
}

func init() {
	rootCmd.AddCommand(genCmd)

}

var genCmd = &cobra.Command{
	Use:   "generate [keyword]",
	Short: "CLI Tool for retrieving SerpAPI results for a Specfied Search Engine",
	Long: `This tool will allow you to get all of the results from a Google search and write it to a CSV file. 
	
	Usage:
	  gdorks generate [Keyword]
	

	Example:
		gdorks generate "John doe"`,
	Run: func(cmd *cobra.Command, args []string) {
		var keyword string
		if len(args) > 0 {
			keyword = args[0]
		} else {
			log.Fatal("Keyword is required. Use --query flag or provide as positional argument")
		}

		for _, val := range allDorks {
			fmt.Println(GenerateDork(string(keyword), val))
		}
	},
}
