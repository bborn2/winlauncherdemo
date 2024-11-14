using System;
using System.Windows.Forms;

namespace MyWindowsApp
{
    public partial class Form1 : Form
    {
        public Form1()
        {
            InitializeComponent();
        }

        private void button1_Click(object sender, EventArgs e)
        {
            MessageBox.Show("This is a dialog!", "Dialog Title", MessageBoxButtons.OK, MessageBoxIcon.Information);
        }
    }
}
