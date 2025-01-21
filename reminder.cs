using System;
using System.Runtime.InteropServices;
 
class Program
{

    [DllImport("reminder.dll", CharSet = CharSet.Unicode, CallingConvention = CallingConvention.Cdecl)]
    public static extern int AddReminder(string  a, string b, string loc, string attends, int dur);
 
    static void Main()
    {
        int result = AddReminder("和YY讨论问题", "2024-11-26 14:00:00", "巴黎会议室", "海总,YY,other", 80);
        Console.WriteLine("Result: " + result);  // 输出结果应该是 8
        Console.ReadLine();
    }
}

// 