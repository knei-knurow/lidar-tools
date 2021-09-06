$output = "output/" + (Get-Date).tostring("dd-MM-yyyy-hh-mm-ss") + ".txt"
$cloudrotation = -3.14159 / 4
./sync.exe `
    --avrport=COM13 `
    --avrbaud=19200 `
    --lidarexe=lidar.exe `
    --lidarport=COM4 `
    --lidarmode=3 `
    --lidarpm=500 `
    --servostep=4 `
    --servodelay=80 `
    --servomin=1000 `
    --servocalib=2500 `
    --servostart=3000 `
    --servomax=3000 `
    --servounit=-0.05 `
    --cloudrotation=$cloudrotation `
    > $output 
