using FaceitReader.Functions;
using System;
using System.Windows.Forms;

namespace FaceitReader.Forms
{
    public partial class Login : Form
    {
        public Login()
        {
            
            InitializeComponent();
            if(Properties.Settings.Default.userName != string.Empty)
            {
                usernameTxt.Text = Properties.Settings.Default.userName;
                passwordTxt.Text = Properties.Settings.Default.passUser;
                if(Properties.Settings.Default.checkBox == true)
                {
                    RememberCheck.Checked = true;
                }
                
            }
        }

        private void button1_Click(object sender, EventArgs e)
        {
            if (!string.IsNullOrWhiteSpace(usernameTxt.Text) && !string.IsNullOrWhiteSpace(passwordTxt.Text))
            {
                if (RememberCheck.Checked)
                {
                    Properties.Settings.Default.userName = usernameTxt.Text;
                    Properties.Settings.Default.passUser = passwordTxt.Text;
                    Properties.Settings.Default.checkBox = true;
                    Properties.Settings.Default.Save();
                }
                else
                {
                    Properties.Settings.Default.userName = "";
                    Properties.Settings.Default.passUser = "";
                    Properties.Settings.Default.checkBox = false;
                    Properties.Settings.Default.Save();
                }
                bool userAuth = Auth.userAuth(usernameTxt.Text, passwordTxt.Text);
                if (userAuth)
                {
                    this.Hide();
                    Main overview = new Main();
                    overview.Show();
                }
                else
                {
                    MessageBox.Show("Username und Passwort stimmen nicht überein!", "Fehler", MessageBoxButtons.OK, MessageBoxIcon.Error);
                }
            }
            else
            {
                MessageBox.Show("Bitte Username und Passwort eingeben!", "Fehler", MessageBoxButtons.OK, MessageBoxIcon.Error);
            }
        }
    }
}