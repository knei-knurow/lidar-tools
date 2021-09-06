$output = "output/" + (Get-Date).tostring("dd-MM-yyyy-hh-mm-ss") + ".txt"
./sync.exe > $output 
