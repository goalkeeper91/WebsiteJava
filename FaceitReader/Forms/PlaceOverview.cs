using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.IO;
using System.Threading.Tasks;
using System.Windows.Forms;

using FaceitReader.Classes;
using FaceitReader.Database;
using FaceitReader.Functions;
using MySql.Data.MySqlClient;
using NPOI.SS.UserModel;
using NPOI.XSSF.UserModel;

namespace FaceitReader
{
    public partial class PlaceOverview : Form
    {
        private String CupId;
        private int Cup;
        private List<FullPlaceList> list;
        private string path;
        private MySqlConnection cnn;
        private MySqlConnection cnn2;

        public PlaceOverview()
        {
            InitializeComponent();
        }

        private void dataGridView1_CellContentClick(object sender, DataGridViewCellEventArgs e)
        {

        }

        private async void send_Click(object sender, EventArgs e)
        {
            int Matches = 0;
            int Places = 0;
            if (!string.IsNullOrWhiteSpace(cupLink.Text))
            {
                if (cupLink.Text.StartsWith("https://www.faceit.com/"))
                {
                    CupId = cupLink.Text;
                    CupId = CupId.Replace("https://www.faceit.com", "");
                    CupId = CupId.Split('/')[3];
                    Cup = cbCup.SelectedIndex;
                }
                else
                {
                    MessageBox.Show("Der eingebene Link ist nicht korrekt!", "Fehler", MessageBoxButtons.OK, MessageBoxIcon.Warning);
                }
            }
            else
            {
                MessageBox.Show("Bitte gib einen Link zum Cup ein.", "Fehler", MessageBoxButtons.OK, MessageBoxIcon.Warning);
            }

            if (Cup == 1)
            {
                Places = 32;
                Matches = 64;
                saveTeamsBtn.Visible = true;
            }
            else if (Cup == 2)
            {
                Places = 8;
                Matches = 31;
                if(saveTeamsBtn.Visible)
                {
                    saveTeamsBtn.Visible = false;
                }
            }
            else
            {
                MessageBox.Show("Du hast keinen Cup ausgewählt", "Fehler", MessageBoxButtons.OK, MessageBoxIcon.Warning);
            }
            list = GetData.getPlacementData(Places,Matches,CupId,Cup);

            var bindingList = new BindingList<FullPlaceList>(list);
            var source = new BindingSource(bindingList, null);

            dataGridView1.DataSource = source;
            dataGridView1.Columns["Id"].Visible = false;
            dataGridView1.Columns["Place"].Width = 45;

            await Task.Run(() => AddValues.fillData(CupId,Matches));
        }

        private void saveTeamsBtn_Click(object sender, EventArgs e)
        {
            var cupName = GetData.getCupName(CupId);
            bool check = false;
            string query = "";
            string queryCheck = "Select * from places where cup_name='"+ cupName +"'";
            MySqlCommand command;
            MySqlDataReader reader1;
            cnn = Connection.openConnection();
            cnn2 = Connection.openConnection();
            IWorkbook workbook = new XSSFWorkbook();
            ISheet sheet1 = workbook.CreateSheet("sheet1");
            command = new MySqlCommand(queryCheck,cnn2);
            reader1 = command.ExecuteReader();
            if (reader1.HasRows)
            {
                check = true;
            }
            reader1.Close();
            cnn2.Close();

            if (folderBrowserDialog1.ShowDialog() == DialogResult.OK)
            {
                path = folderBrowserDialog1.SelectedPath;
                var x = 2;
                var Points = 0;
                IRow row0 = sheet1.CreateRow(0);
                row0.CreateCell(0).SetCellValue(cupName);
                IRow row1 = sheet1.CreateRow(1);
                row1.CreateCell(0).SetCellValue("Platz");
                row1.CreateCell(1).SetCellValue("Teamname");
                row1.CreateCell(2).SetCellValue("Teamcaptain");
                row1.CreateCell(3).SetCellValue("Punkte");

                using (var file = File.Create(path+@"\teams_"+CupId+".xlsx"))
                {
                    foreach (var arr in list)
                    {
                        if(arr.Place <= 8)
                        {
                            if(arr.Place == 1)
                            {
                                Points = 25;
                            }
                            else if(arr.Place == 2)
                            {
                                Points = 20;
                            }
                            else if(arr.Place == 3 || arr.Place == 4)
                            {
                                Points = 15;
                            }
                            else
                            {
                                Points = 10;
                            }
                            if (!check)
                            {
                                query = "INSERT INTO places (teamname,points,cup_name) VALUES('" + arr.TeamName + "','" + Points + "','" + cupName + "')";
                                command = new MySqlCommand(query,cnn);
                                command.ExecuteNonQuery();
                            }
                            IRow row2  = sheet1.CreateRow(x);
                            row2.CreateCell(0).SetCellValue(arr.Place);
                            row2.CreateCell(1).SetCellValue(arr.TeamName);
                            row2.CreateCell(2).SetCellValue(arr.CaptainName);
                            row2.CreateCell(3).SetCellValue(Points);
                            x++;
                        }
                    }
                    cnn.Close();
                    workbook.Write(file);
                    file.Close();
                }
            }
            
        }

        private void PlaceOverview_FormClosed(object sender, FormClosedEventArgs e)
        {
            Application.Exit();
        }

        private void btnClose_Click(object sender, EventArgs e)
        {
            this.Close();
        }
    }
}
