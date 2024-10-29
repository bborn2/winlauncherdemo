// mac.go
//go:build darwin
// +build darwin

package main

func GetInstalledPrograms() ([]Program, error) {

	return []Program{
		{Name: "Git", Path: `C:\Program Files\Git\mingw64\share\git\git-for-windows.ico`},
		{Name: "Lenovo AI Now 1.0", Path: `C:\Program Files\Lenovo\Lenovo AI Now\Lenovo AINow.exe`},
		{Name: "Logitech Unifying ¨¨¨ª?t 2.52", Path: `C:\Program Files\Common Files\LogiShrd\Unifying\DJCUHost.exe`},
		{Name: "Microsoft 365 - en-us", Path: `C:\Program Files\Common Files\Microsoft Shared\ClickToRun\OfficeClickToRun.exe`},
		{Name: "Symantec Endpoint Protection", Path: `C:\Program Files\Symantec\Symantec Endpoint Protection\14.3.9681.7000.105\Bin64\SymCorpUI.exe`},
		{Name: "Malwarebytes version 4.6.21.347", Path: `C:\Program Files\Malwarebytes\Anti-Malware\mbam.exe`},
		{Name: "Cisco AnyConnect Secure Mobility Client", Path: `C:\Program Files (x86)\Cisco\Cisco AnyConnect Secure Mobility Client\InstallHelper.exe`},
		{Name: "Google Chrome", Path: `C:\Program Files\Google\Chrome\Application\chrome.exe,0`},
		{Name: "Microsoft Edge", Path: `C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe,0`},
		{Name: "Microsoft Edge WebView2 Runtime", Path: `C:\Program Files (x86)\Microsoft\EdgeWebView\Application\130.0.2849.52\msedgewebview2.exe,0`},
		{Name: "PyCharm Community Edition 2024.2", Path: `C:\Program Files\JetBrains\PyCharm Community Edition 2024.2\bin\pycharm64.exe`},
		{Name: "Lenovo Vantage Service", Path: `C:\Program Files (x86)\Lenovo\VantageService\4.2.24.0\\Uninstall.exe`},
		{Name: "微信", Path: `"C:\Program Files\Tencent\WeChat\WeChat.exe"`},
		{Name: "Cybereason Sensor", Path: `C:\ProgramData\Package Cache\{25f00eb4-883a-434f-a22d-ca17d4706190}\CybereasonSensor.exe,0`},
		{Name: "Cisco AnyConnect ISE Posture Module", Path: `C:\Program Files (x86)\Cisco\Cisco AnyConnect Secure Mobility Client\InstallHelper.exe`},
		{Name: "Malwarebytes Endpoint Agent", Path: `C:\ProgramData\Package Cache\{acd88974-dfd2-4107-a1b0-ba3be6b6199e}\Setup.MBEndpointAgent.Full.exe,0`},
		{Name: "Lenovo Quick Clean", Path: `C:\Program Files (x86)\Lenovo\Lenovo Quick Clean\LenovoQuickClean.exe`},
		{Name: "百度网盘", Path: `"C:\Users\T14S\AppData\Roaming\baidu\BaiduNetdisk\BaiduNetdisk.exe"`},
	}, nil
}
