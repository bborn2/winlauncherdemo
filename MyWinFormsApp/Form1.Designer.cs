﻿namespace MyWinFormsApp;

partial class Form1
{
    /// <summary>
    ///  Required designer variable.
    /// </summary>
    private System.ComponentModel.IContainer components = null;

    /// <summary>
    ///  Clean up any resources being used.
    /// </summary>
    /// <param name="disposing">true if managed resources should be disposed; otherwise, false.</param>
    protected override void Dispose(bool disposing)
    {
        if (disposing && (components != null))
        {
            components.Dispose();
        }
        base.Dispose(disposing);
    }

    #region Windows Form Designer generated code

    /// <summary>
    ///  Required method for Designer support - do not modify
    ///  the contents of this method with the code editor.
    /// </summary>
    // private void InitializeComponent()
    // {
    //     this.components = new System.ComponentModel.Container();
    //     this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
    //     this.ClientSize = new System.Drawing.Size(800, 450);
    //     this.Text = "Form1";
    // }

            private void InitializeComponent()
        {
            this.button1 = new Button();
            this.SuspendLayout();

            // 
            // button1
            // 
            this.button1.Location = new System.Drawing.Point(100, 50);
            this.button1.Name = "button1";
            this.button1.Size = new System.Drawing.Size(100, 50);
            this.button1.Text = "Click Me";
            this.button1.UseVisualStyleBackColor = true;
            this.button1.Click += new EventHandler(this.button1_Click);

            // 
            // Form1
            // 
            this.ClientSize = new System.Drawing.Size(300, 150);
            this.Controls.Add(this.button1);
            this.Name = "Form1";
            this.Text = "My WinForms App";
            this.ResumeLayout(false);
        }

    #endregion
}