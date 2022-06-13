using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Forms;

namespace FaceitReader.Forms
{
    public partial class Main : Form
    {
        public Main()
        {
            InitializeComponent();
            if (Globals.admin)
            {
                btnAdmin.Show();
            }
        }

        private void btnAdmin_Click(object sender, EventArgs e)
        {

        }
        private void btnEndCup_Click(object sender, EventArgs e)
        {
            PlaceOverview endCup = new PlaceOverview();
            endCup.Show();
        }

        private void btnBanTeam_Click(object sender, EventArgs e)
        {

        }

        private void btnTeamCheck_Click(object sender, EventArgs e)
        {
            TeamCheckModal checkTeam = new TeamCheckModal();
            checkTeam.Show();
        }
        private void btnLogout_Click(object sender, EventArgs e)
        {
            Login login = new Login();
            login.Show();
            this.Close();
        }

        
    }
}
