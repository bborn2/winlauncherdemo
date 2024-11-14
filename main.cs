using System;
using System.Runtime.InteropServices;
 
class Program
{

    [DllImport("fastlaunch.dll", CallingConvention = CallingConvention.Cdecl)]
    // public static extern int Add(int a, int b);
    // public static extern int AddReminder(string  a, string b);
    // public static extern void PrintMessage(string msg);
    public static extern void loadApps();

    [DllImport("fastlaunch.dll", CallingConvention = CallingConvention.Cdecl)]
    public static extern int searchAndRun(string query);

    static void Main()
    {
        loadApps();

        int result = searchAndRun("word");
        // int result = Add(1,2);
        // Console.WriteLine("Result: " + result);  // 输出结果应该是 8

        // PrintMessage("hello world");

        // int result = AddReminder("test", "2024-11-09 14:00:00");
        Console.WriteLine("Result: " + result);  // 输出结果应该是 8
        Console.ReadLine();
    }
}

// C:\Windows\Microsoft.NET\Framework64\v3.5\csc.exe .\main.cs