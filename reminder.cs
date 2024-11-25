using System;
using System.Runtime.InteropServices;
 
class Program
{

    [DllImport("reminder.dll", CharSet = CharSet.Unicode, CallingConvention = CallingConvention.Cdecl)]
    public static extern int AddReminder(string  a, string b);
 
    static void Main()
    {
        int result = AddReminder("和YY讨论问题", "2024-11-26 14:00:00");
        Console.WriteLine("Result: " + result);  // 输出结果应该是 8
        Console.ReadLine();
    }
}

// 