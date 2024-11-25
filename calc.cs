using System;
using Outlook = Microsoft.Office.Interop.Outlook;

namespace OutlookMeeting
{
    class Program
    {
        static void Main(string[] args)
        {
            Outlook.Application outlookApp = new Outlook.Application();
            Outlook.Namespace outlookNamespace = outlookApp.GetNamespace("MAPI");
            Outlook.AppointmentItem appointmentItem = outlookNamespace.CreateItem(Outlook.OlItemType.olAppointmentItem);

            appointmentItem.Subject = "新会议邀请";
            // ... 其他属性设置

            appointmentItem.Display(true);
        }
    }
}