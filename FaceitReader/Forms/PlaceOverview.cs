using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.IO;
using System.Threading;
using System.Windows.Forms;
using FaceitReader.Classes;
using FaceitReader.Functions;
using Google.Apis.Auth.OAuth2;
using Google.Apis.Services;
using Google.Apis.Sheets.v4.Data;
using Google.Apis.Sheets.v4;
using Google.Apis.Util.Store;
using Color = System.Drawing.Color;
using System.Linq.Expressions;
using System.Linq;

namespace FaceitReader
{
    public partial class PlaceOverview : Form
    {
        private string CupId;
        private int Cup;
        string Cupname;
        private List<FullPlaceList> list;
        static readonly string[] Scopes = { SheetsService.Scope.Spreadsheets };
        static readonly string ApplicationName = "Faceit Reader";
        string range;
        string voucherRange;
        string sheetName;
        readonly string version = Application.ProductVersion.ToString();
        

        // Define request parameters.
        string spreadsheetId;

        public PlaceOverview()
        {
            InitializeComponent();
            this.txtSearch.Enter += new EventHandler(txtSearch_Enter);
            this.txtSearch.Leave += new EventHandler(txtSearch_Leave);
            txtSearch_SetText();
            dataGridView1.LostFocus += Grid_LostFocus;
            dataGridView2.LostFocus += Grid_LostFocus;
            toolStripStatusLabel2.Text = version;
        }
        private void Grid_LostFocus(object sender, EventArgs e)
        {
            dataGridView1.ClearSelection();
            dataGridView2.ClearSelection();
        }
        protected void txtSearch_SetText()
        {
            txtSearch.Text = "Teamname oder Captain";
            txtSearch.ForeColor = Color.Gray;
        }
        private void txtSearch_Enter(object sender, EventArgs e)
        {
            if (txtSearch.ForeColor == Color.Black)
                return;
            txtSearch.Text = "";
            txtSearch.ForeColor = Color.Black;
        }
        private void txtSearch_Leave(object sender, EventArgs e)
        {
            if (txtSearch.Text.Trim() == "")
                txtSearch_SetText();
        }
        private void send_Click(object sender, EventArgs e)
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
                    Cupname = cupLink.Text.Replace("https://www.faceit.com", "").Split('/')[4].Replace("%20", " ");
                    Cup = cbCup.SelectedIndex;
                    setRange(Cupname, Cup);
                    setVoucherRange(Cupname, Cup);
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

            if(Cupname.Contains("Polska") && Cup == 1)
            { 
                DialogResult dialogResult = MessageBox.Show("Der Cuplink und der ausgewählte Cup stimmen nicht überein! Möchtest du trotzdem Fortfahren?", "Achtung", MessageBoxButtons.YesNoCancel, MessageBoxIcon.Warning); 
                if( dialogResult == DialogResult.Yes)
                {
                    if (Cup == 1)
                    {
                        Places = 32;
                        Matches = 64;
                        spreadsheetId = "1hZZsHOKiEaVUy7tTCnMtjuTkh0rzpsHhr1guYHFLgzo";
                    }
                    else if (Cup == 2)
                    {
                        Places = 8;
                        Matches = 31;
                        spreadsheetId = "1mqrtmvGvcfh1p6slI8PR8Ya-j8Jizw-iyNI4y96dtvE";
                    }
                }
                else
                {
                    cupLink.Text = string.Empty;
                    cbCup.SelectedIndex = 0;
                    return;
                }
            }
            else if (!Cupname.Contains("Polska") && Cup == 2)
            {
                DialogResult dialogResult = MessageBox.Show("Der Cuplink und der ausgewählte Cup stimmen nicht überein! Möchtest du trotzdem Fortfahren?", "Achtung", MessageBoxButtons.YesNoCancel, MessageBoxIcon.Warning);
                if (dialogResult == DialogResult.Yes)
                {
                    if (Cup == 1)
                    {
                        Places = 32;
                        Matches = 64;
                        spreadsheetId = "1hZZsHOKiEaVUy7tTCnMtjuTkh0rzpsHhr1guYHFLgzo";
                    }
                    else if (Cup == 2)
                    {
                        Places = 8;
                        Matches = 31;
                        spreadsheetId = "1mqrtmvGvcfh1p6slI8PR8Ya-j8Jizw-iyNI4y96dtvE";
                    }
                }
                else
                {
                    cupLink.Text = string.Empty;
                    cbCup.SelectedIndex = 0;
                    return;
                }
            }
            else if (Cup == 1 || Cup == 2)
            {
                if (Cup == 1)
                {
                    Places = 32;
                    Matches = 64;
                    spreadsheetId = "1hZZsHOKiEaVUy7tTCnMtjuTkh0rzpsHhr1guYHFLgzo";
                }
                else if (Cup == 2)
                {
                    Places = 8;
                    Matches = 31;
                    spreadsheetId = "1mqrtmvGvcfh1p6slI8PR8Ya-j8Jizw-iyNI4y96dtvE";
                }
            }
            else
            {
                MessageBox.Show("Du hast keinen Cup ausgewählt", "Fehler", MessageBoxButtons.OK, MessageBoxIcon.Warning);
            }

            list = GetData.getPlacementData(Places, Matches, CupId, Cup);

            var bindingList = new BindingList<FullPlaceList>(list);
            var source = new BindingSource(bindingList, null);

            dataGridView1.DataSource = source;
            dataGridView1.Columns["Id"].Visible = false;
            dataGridView1.Columns["TeamId"].Visible = false;
            dataGridView1.Columns["Place"].Width = 45;
            dataGridView1.CellFormatting += new DataGridViewCellFormattingEventHandler(this.dataGridView1_CellFormatting);

            try 
            {
                if (Cup == 1 || Cup == 2)
                {
                    UserCredential credential;

                    using (var stream =
                        new FileStream("Resources/credentials.json", FileMode.Open, FileAccess.Read))
                    {
                        string credPath = "token.json";
                        credential = GoogleWebAuthorizationBroker.AuthorizeAsync(
                            GoogleClientSecrets.FromStream(stream).Secrets,
                            Scopes,
                            "user",
                            CancellationToken.None,
                            new FileDataStore(credPath, true)).Result;
                    }
                    var service = new SheetsService(new BaseClientService.Initializer()
                    {
                        HttpClientInitializer = credential,
                        ApplicationName = ApplicationName,
                    });

                    SpreadsheetsResource.ValuesResource.GetRequest request =
                            service.Spreadsheets.Values.Get(spreadsheetId, voucherRange);
                    SpreadsheetsResource.GetRequest request2 = service.Spreadsheets.Get(spreadsheetId);
                    request2.Ranges = range;
                    request2.IncludeGridData = true;
                    // https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit

                    ValueRange response = request.Execute();
                    IList<IList<Object>> values = response.Values;
                    if (values != null && values.Count > 0)
                    {
                        var bindinglist3 = new BindingList<Voucher>(VoucherMapper.MapFromRangeData(values, Cup));

                        dataGridView2.DataSource = bindinglist3;
                        dataGridView2.Columns["Id"].Visible = false;
                        if (Cup == 1)
                        {
                            dataGridView2.Columns[1].HeaderText = "25€";
                            dataGridView2.Columns[2].HeaderText = "15€";
                            dataGridView2.Columns[3].HeaderText = "10€";
                            dataGridView2.Columns[4].HeaderText = "5€";
                            dataGridView2.Columns[5].HeaderText = "2€";
                        }
                        else if (Cup == 2)
                        {
                            dataGridView2.Columns[1].HeaderText = "25€";
                            dataGridView2.Columns[2].Visible = false;
                            dataGridView2.Columns[3].HeaderText = "10€";
                            dataGridView2.Columns[4].HeaderText = "5€";
                            dataGridView2.Columns[5].HeaderText = "2€";
                        }
                        dataGridView2.CellFormatting += new DataGridViewCellFormattingEventHandler(this.dataGridView2_CellFormatting);
                    }
                    else
                    {
                        MessageBox.Show("Keine Daten gefunden", "Error", MessageBoxButtons.OK, MessageBoxIcon.Error);
                    }
                    Spreadsheet response2 = request2.Execute();
                }
            }
            catch (Exception ex)
            {
                dataGridView2.DataSource = null;
                string message;
                if (ex.Message.Contains("range"))
                {
                    message = "Für den Cup konnten keine Voucher gefunden werden. Entweder verwendest du einen falschen Link oder der Cup ist aus einer früheren Saison!";
                }
                else
                {
                    message = ex.Message;
                }
                MessageBox.Show(message, "Error", MessageBoxButtons.OK, MessageBoxIcon.Error);
            }
            
        }
        private void dataGridView1_CellFormatting(object sender, DataGridViewCellFormattingEventArgs e)
        {
            foreach (DataGridViewRow row in dataGridView1.Rows)
            {
                if (row.Cells[2].Value.ToString().Contains("team_"))
                {
                    row.DefaultCellStyle.BackColor = Color.DarkSalmon;
                }
            }
        }
        private void dataGridView2_CellFormatting(object sender, DataGridViewCellFormattingEventArgs e)
        {
            int i = 4;
            int j = 10;
            dataGridView2.RowsDefaultCellStyle.SelectionBackColor = Color.ForestGreen;
            foreach ( DataGridViewRow row in dataGridView2.Rows)
            {
                if (row.Index > i && row.Index < j)
                {
                    row.DefaultCellStyle.BackColor = Color.Gray;
                }
                else
                {
                    row.DefaultCellStyle.BackColor = Color.DodgerBlue;
                }
                if (row.Index == j)
                {
                    i += 10;
                    j += 10;
                }
            }
        }
        private void btnClose_Click(object sender, EventArgs e)
        {
            this.Close();
        }
        private void btnSearch_Click(object sender, EventArgs e)
        {
            string searchValue = txtSearch.Text.ToLower();
            dataGridView1.ClearSelection();

            dataGridView1.SelectionMode = DataGridViewSelectionMode.FullRowSelect;
            try
            {
                foreach (DataGridViewRow row in dataGridView1.Rows)
                {
                    if (row.Cells[2].Value.ToString().ToLower().Contains(searchValue))
                    {
                        row.Selected = true;
                        dataGridView1.FirstDisplayedScrollingRowIndex = dataGridView1.SelectedRows[0].Index;
                        break;
                    }
                    else if (row.Cells[4].Value.ToString().ToLower().Contains(searchValue))
                    {
                        row.Selected = true;
                        dataGridView1.FirstDisplayedScrollingRowIndex = dataGridView1.SelectedRows[0].Index;
                        break;
                    }
                    else
                    {
                        MessageBox.Show("Kein Team oder Captain stimmt mit der Eingabe überein!", "Achtung", MessageBoxButtons.OK, MessageBoxIcon.Information);
                        break;
                    }
                }
            }
            catch (Exception exc)
            {
                MessageBox.Show(exc.Message);
            }
        }
        private void btn_save_Click(object sender, EventArgs e)
        {
            try
            {
                GetData.CreateEntry(dataGridView1, spreadsheetId, range);
            }
            catch (Exception ex)
            {
                MessageBox.Show(ex.Message.ToString(), "Error", MessageBoxButtons.OK, MessageBoxIcon.Error);
            }
        }
        public void setRange(string name, int id)
        {
            if(id == 1)
            {
                if (name.Contains("Cup 10"))
                {
                    range = "Übersicht!J45:J76";
                }
                else if (name.Contains("Cup 2"))
                {
                    range = "Übersicht!D7:D38";
                }
                else if (name.Contains("Cup 3"))
                {
                    range = "Übersicht!F7:F38";
                }
                else if (name.Contains("Cup 4"))
                {
                    range = "Übersicht!H7:H38";
                }
                else if (name.Contains("Cup 5"))
                {
                    range = "Übersicht!J7:J38";
                }
                else if (name.Contains("Cup 6"))
                {
                    range = "Übersicht!B45:B76";
                }
                else if (name.Contains("Cup 7"))
                {
                    range = "Übersicht!D45:D76";
                }
                else if (name.Contains("Cup 8"))
                {
                    range = "Übersicht!F45:F76";
                }
                else if (name.Contains("Cup 9"))
                {
                    range = "Übersicht!H45:H76";
                }
                else if (name.Contains("Cup 1"))
                {
                    range = "Übersicht!B7:B38";
                }
            }
            else if(id == 2)
            {
                if (name.Contains("C10"))
                {
                    range = "Übersicht!J20:J27";
                }
                else if (name.Contains("C2"))
                {
                    range = "Übersicht!D7:D14";
                }
                else if (name.Contains("C3"))
                {
                    range = "Übersicht!F7:F14";
                }
                else if (name.Contains("C4"))
                {
                    range = "Übersicht!H7:H14";
                }
                else if (name.Contains("C5"))
                {
                    range = "Übersicht!J7:J14";
                }
                else if (name.Contains("C6"))
                {
                    range = "Übersicht!B20:B27";
                }
                else if (name.Contains("C7"))
                {
                    range = "Übersicht!D20:D27";
                }
                else if (name.Contains("C8"))
                {
                    range = "Übersicht!F20:F27";
                }
                else if (name.Contains("C9"))
                {
                    range = "Übersicht!H20:H27";
                }
                else if (name.Contains("C1"))
                {
                    range = "Übersicht!B7:B14";
                }
            }
        }
        public void setVoucherRange(string name, int id)
        {
            string r = "!A2:E121";
            string p = "!A2:D21";
            if (id == 1)
            {
                if (name.Contains("Cup 10"))
                {
                    voucherRange = Cupname.Replace("SKINBARON CUP ", "") + r;
                }
                else if (name.Contains("Cup 2"))
                {
                    voucherRange = Cupname.Replace("SKINBARON CUP ", "") + r;
                }
                else if (name.Contains("Cup 3"))
                {
                    voucherRange = Cupname.Replace("SKINBARON CUP ", "") + r;
                }
                else if (name.Contains("Cup 4"))
                {
                    voucherRange = Cupname.Replace("SKINBARON CUP ", "") + r;
                }
                else if (name.Contains("Cup 5"))
                {
                    voucherRange = Cupname.Replace("SKINBARON CUP ", "") + r;
                }
                else if (name.Contains("Cup 6"))
                {
                    voucherRange = Cupname.Replace("SKINBARON CUP ", "") + r;
                }
                else if (name.Contains("Cup 7"))
                {
                    voucherRange = Cupname.Replace("SKINBARON CUP ", "") + r;
                }
                else if (name.Contains("Cup 8"))
                {
                    voucherRange = Cupname.Replace("SKINBARON CUP ", "") + r;
                }
                else if (name.Contains("Cup 9"))
                {
                    voucherRange = Cupname.Replace("SKINBARON CUP ", "") + r;
                }
                else if (name.Contains("Cup 1"))
                {
                    voucherRange = Cupname.Replace("SKINBARON CUP ", "") + r;
                }
                sheetName = Cupname.Replace("SKINBARON CUP ", "");
            }
            else if (id == 2)
            {
                if (name.Contains("C10"))
                {
                    voucherRange = "PL Cup 10" + p;
                    sheetName = "PL Cup 10";
                }
                else if (name.Contains("C2"))
                {
                    voucherRange = "PL Cup 2" + p;
                    sheetName = "PL Cup 2";
                }
                else if (name.Contains("C3"))
                {
                    voucherRange = "PL Cup 3" + p;
                    sheetName = "PL Cup 3";
                }
                else if (name.Contains("C4"))
                {
                    voucherRange = "PL Cup 4" + p;
                    sheetName = "PL Cup 4";
                }
                else if (name.Contains("C5"))
                {
                    voucherRange = "PL Cup 5" + p;
                    sheetName = "PL Cup 5";
                }
                else if (name.Contains("C6"))
                {
                    voucherRange = "PL Cup 6" + p;
                    sheetName = "PL Cup 6";
                }
                else if (name.Contains("C7"))
                {
                    voucherRange = "PL Cup 7" + p;
                    sheetName = "PL Cup 7";
                }
                else if (name.Contains("C8"))
                {
                    voucherRange = "PL Cup 8" + p;
                    sheetName = "PL Cup 8";
                }
                else if (name.Contains("C9"))
                {
                    voucherRange = "PL Cup 9" + p;
                    sheetName = "PL Cup 9";
                }
                else if (name.Contains("C1"))
                {
                    voucherRange = "PL Cup 1" + p;
                    sheetName = "PL Cup 1";
                }
            }
        }
        private void txtSearch_KeyDown(object sender, KeyEventArgs e)
        {
            if (e.KeyCode == Keys.Enter)
            {
                string searchValue = txtSearch.Text;
                dataGridView1.ClearSelection();

                dataGridView1.SelectionMode = DataGridViewSelectionMode.FullRowSelect;
                try
                {
                    foreach (DataGridViewRow row in dataGridView1.Rows)
                    {
                        if (row.Cells[2].Value.ToString().ToLower().Contains(searchValue))
                        {
                            row.Selected = true;
                            dataGridView1.FirstDisplayedScrollingRowIndex = dataGridView1.SelectedRows[0].Index;
                            break;
                        }
                        else if (row.Cells[4].Value.ToString().ToLower().Contains(searchValue))
                        {
                            row.Selected = true;
                            dataGridView1.FirstDisplayedScrollingRowIndex = dataGridView1.SelectedRows[0].Index;
                            break;
                        }
                    }
                    e.Handled = true;
                    e.SuppressKeyPress = true;
                }
                catch (Exception exc)
                {
                    MessageBox.Show(exc.Message);
                }
            }
        }
        private void Lost_Focus(object sender, MouseEventArgs e)
        {
            Grid_LostFocus(sender, e);
        }
        private void dataGridView1_CellContentDoubleClick(object sender, DataGridViewCellEventArgs e)
        {
            int[] value = setRow(e.RowIndex);
            if(value[1] != 0) 
            {
                DialogResult result = MessageBox.Show("Hast du die Codes wirklich an " + dataGridView1.Rows[e.RowIndex].Cells[2].Value.ToString() + " gesendet?", "Warning", MessageBoxButtons.YesNo, MessageBoxIcon.Question);
                if (result == DialogResult.Yes)
                {
                    try
                    {
                        GetData.updateSpreadsheetCells(spreadsheetId, "Übersicht", value[0], value[1], value[1], dataGridView1.Rows[e.RowIndex].Cells[2].Value.ToString(), dataGridView1.Rows[e.RowIndex].Cells[4].Value.ToString());
                        MessageBox.Show("Das Team "+ dataGridView1.Rows[e.RowIndex].Cells[2].Value.ToString()+" wurde heute" + DateTime.Now.ToString("dd.MM.yyyy") + " als, Voucher erhalten, markiert!", "Bestätigung", MessageBoxButtons.OK, MessageBoxIcon.Information);
                    }
                    catch ( Exception exc )
                    {
                        MessageBox.Show(exc.ToString(), "Error", MessageBoxButtons.OK, MessageBoxIcon.Error);
                    }
                }
            }
        }
        private int[] setRow(int index)
        {
            int col;
            int row;
            int[] ints = new int[2];
            if (Cupname.Contains("1") || Cupname.Contains("6"))
            {
                col = 1;
                ints[0] = col;
            }
            else if (Cupname.Contains("2") || Cupname.Contains("7"))
            {
                col = 3;
                ints[0] = col;
            }
            else if (Cupname.Contains("3") || Cupname.Contains("8"))
            {
                col = 5;
                ints[0] = col;
            }
            else if (Cupname.Contains("4") || Cupname.Contains("9"))
            {
                col = 7;
                ints[0] = col;
            }
            else if (Cupname.Contains("5") || Cupname.Contains("10"))
            {
                col = 9;
                ints[0] = col;
            }
            if (Cupname.Contains("1") || Cupname.Contains("2") || Cupname.Contains("3")  || Cupname.Contains("4")  || Cupname.Contains("5"))
            {
                row = index + 6;
                ints[1] = row;
                return ints;
            }
            else if(Cupname.Contains("6") || Cupname.Contains("7") || Cupname.Contains("8") || Cupname.Contains("9") || Cupname.Contains("10"))
            {
                row = index + 44;
                ints[1] = row;
                return ints;
            }

            return ints;
        }
        private void dataGridView2_MouseDown(object sender, MouseEventArgs e)
        {
            List<int> cols =  new List<int>();
            int y = 0;
            int x = 0;
            if(e.Button == MouseButtons.Right) 
            {
                int selectedCellCount = dataGridView2.SelectedCells.Count;
                List<int> rows = new List<int>();
                for (int i = 0; i < selectedCellCount; i++)
                {
                    y = dataGridView2.SelectedCells[i].ColumnIndex;
                    if(y != x)
                    {
                        cols.Add(dataGridView2.SelectedCells[i].ColumnIndex);
                        x = y;
                    }
                    rows.Add(dataGridView2.SelectedCells[i].RowIndex);
                }
                if(cols.Count > 1) 
                {
                    MessageBox.Show("Bitte wähle nur eine Spalte zur Zeit", "Warnung", MessageBoxButtons.OK, MessageBoxIcon.Warning);
                }
                else
                {
                    try
                    {
                        if(selectedCellCount % 5 == 0)
                        {
                            DialogResult result = MessageBox.Show("Möchtest du die Voucher als versendet markieren?", "Voucher markieren?", MessageBoxButtons.YesNo, MessageBoxIcon.Question);
                            if(result == DialogResult.Yes)
                            {
                                if (cols[0] - 1 == 0)
                                {
                                    GetData.updateSpreadsheetCells(spreadsheetId, sheetName, cols[0] - 1, rows.Last() + 1, rows.First() + 1, null, null);
                                    MessageBox.Show("Deine ausgewählten Voucher wurden als versendet Markiert", "Bestätigung", MessageBoxButtons.OK, MessageBoxIcon.Information);
                                }
                                else
                                {
                                    GetData.updateSpreadsheetCells(spreadsheetId, sheetName, cols[0] - 2, rows.Last() + 1, rows.First() + 1, null, null);
                                    MessageBox.Show("Deine ausgewählten Voucher wurden als versendet Markiert", "Bestätigung", MessageBoxButtons.OK, MessageBoxIcon.Information);
                                }
                            }
                        }
                        else
                        {
                            DialogResult result = MessageBox.Show("Du hast "+ selectedCellCount +" ausgewählt. Bist du dir sicher, dass du die Voucher verschickt hast?", "Überprüfe deine Auswahl", MessageBoxButtons.YesNoCancel, MessageBoxIcon.Question);
                            if(result == DialogResult.Yes)
                            {
                                GetData.updateSpreadsheetCells(spreadsheetId, sheetName, cols[0] - 1, rows.Last() + 1, rows.First() + 1, null, null);
                                MessageBox.Show("Deine ausgewählten Voucher wurden als versendet Markiert", "Bestätigung", MessageBoxButtons.OK, MessageBoxIcon.Information);
                            }
                        }
                    }
                    catch (Exception exc)
                    {
                        MessageBox.Show(exc.ToString(), "Error", MessageBoxButtons.OK, MessageBoxIcon.Error);
                    }
                }
            }
        }
    }
}
