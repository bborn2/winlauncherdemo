using System;
using System.Runtime.InteropServices;
 
class Program
{

    [DllImport("reminder.dll", CallingConvention = CallingConvention.Cdecl)]
    public static extern int AddReminder(string  a, string b);
 
    static void Main()
    {
        int result = AddReminder("test", "2024-11-16 14:00:00");
        Console.WriteLine("Result: " + result);  // 输出结果应该是 8
        Console.ReadLine();
    }
}

// C:\Windows\Microsoft.NET\Framework64\v3.5\csc.exe .\main.cs