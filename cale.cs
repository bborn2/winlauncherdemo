using System;
using System.Runtime.InteropServices;
 
class Program
{

    [DllImport("cale.dll", CharSet = CharSet.Unicode, CallingConvention = CallingConvention.Cdecl)]
    public static extern int AddCale(string  topic, string date, string attendees);
  
    static void Main()
    {
        int result = AddCale("meeting invitation", "2025-03-10", "shuangyang;song kun;lihua");
        Console.WriteLine("Result: " + result);
        Console.ReadLine();
    }
}

// C:\Windows\Microsoft.NET\Framework64\v3.5\csc.exe .\cale.cs