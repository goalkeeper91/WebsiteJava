using FaceitReader.Functions;
using System;
using System.Collections.Generic;
using System.Windows.Forms;

namespace FaceitReader.Forms
{
    public partial class TeamCheckModal : Form
    {
        private int limit;
        List<string> teams = new List<string>();
        private string CupId;
        public TeamCheckModal()
        {
            InitializeComponent();
        }

        private void btnOk_Click(object sender, EventArgs e)
        {
            if (cmbCup.SelectedIndex != 0 && !string.IsNullOrWhiteSpace(txtLink.Text))
            {
                if (txtLink.Text.StartsWith("https://www.faceit.com/"))
                {
                    CupId = txtLink.Text;
                    CupId = CupId.Replace("https://www.faceit.com", "");
                    CupId = CupId.Split('/')[3];
                }
                else
                {
                    MessageBox.Show("Der eingebene Link ist nicht korrekt!", "Fehler", MessageBoxButtons.OK, MessageBoxIcon.Warning);
                }
                if (cmbCup.SelectedIndex == 1)
                {
                    limit = 32;
                }
                else if (cmbCup.SelectedIndex == 2)
                {
                    limit = 128;
                }
                teams = CheckTeam.isTeamBanned(txtLink.Text, limit);

                if (teams.Count > 0)
                {
                    string t = "";
                    foreach (string team in teams)
                    {
                        t = t + ", " + team;
                    }
                    MessageBox.Show("Die folgenden Teams sind bei unseren Cups nicht spielberechtigt: " + t, "Info", MessageBoxButtons.OK, MessageBoxIcon.Warning);
                }
                else
                {
                    MessageBox.Show("Alles in Ordnung", "Info", MessageBoxButtons.OK, MessageBoxIcon.Information);
                }
            }
            else
            {
                MessageBox.Show("Bitte fülle die Felder richtig aus!", "Fehler", MessageBoxButtons.OK, MessageBoxIcon.Warning);
            }
        }

        private void btnCancel_Click(object sender, EventArgs e)
        {
            this.Close();
        }
    }
}
