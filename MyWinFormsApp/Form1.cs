using System;
using System.Windows.Forms;
using System.Runtime.InteropServices;

namespace MyWinFormsApp
{    
    public partial class Form1 : Form
    {

            [DllImport("fastlaunch.dll", CallingConvention = CallingConvention.Cdecl)]
    public static extern void loadApps();

    [DllImport("fastlaunch.dll", CallingConvention = CallingConvention.Cdecl)]
    public static extern int searchAndRun(string query);


        private Button button1;

        public Form1()
        {
            InitializeComponent();
 
            loadApps();
        }



        private void button1_Click(object sender, EventArgs e)
        {
            // MessageBox.Show("This is a dialog!", "Dialog Title", MessageBoxButtons.OK, MessageBoxIcon.Information);

            searchAndRun("huatu");
        }
    }
}
