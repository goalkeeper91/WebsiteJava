namespace FaceitReader.Forms
{
    partial class Main
    {
        /// <summary>
        /// Required designer variable.
        /// </summary>
        private System.ComponentModel.IContainer components = null;

        /// <summary>
        /// Clean up any resources being used.
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
        /// Required method for Designer support - do not modify
        /// the contents of this method with the code editor.
        /// </summary>
        private void InitializeComponent()
        {
            System.ComponentModel.ComponentResourceManager resources = new System.ComponentModel.ComponentResourceManager(typeof(Main));
            this.groupBoxHead = new System.Windows.Forms.GroupBox();
            this.btnAdmin = new System.Windows.Forms.Button();
            this.btnBanTeam = new System.Windows.Forms.Button();
            this.btnTeamCheck = new System.Windows.Forms.Button();
            this.btnEndCup = new System.Windows.Forms.Button();
            this.btnLogout = new System.Windows.Forms.Button();
            this.groupBoxValue = new System.Windows.Forms.GroupBox();
            this.tabControl1 = new System.Windows.Forms.TabControl();
            this.tabPage1 = new System.Windows.Forms.TabPage();
            this.dtGrVPolska = new System.Windows.Forms.DataGridView();
            this.tabPage2 = new System.Windows.Forms.TabPage();
            this.dtGrVSkBCup = new System.Windows.Forms.DataGridView();
            this.groupBoxHead.SuspendLayout();
            this.groupBoxValue.SuspendLayout();
            this.tabControl1.SuspendLayout();
            this.tabPage1.SuspendLayout();
            ((System.ComponentModel.ISupportInitialize)(this.dtGrVPolska)).BeginInit();
            this.tabPage2.SuspendLayout();
            ((System.ComponentModel.ISupportInitialize)(this.dtGrVSkBCup)).BeginInit();
            this.SuspendLayout();
            // 
            // groupBoxHead
            // 
            this.groupBoxHead.Anchor = ((System.Windows.Forms.AnchorStyles)((((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Bottom) 
            | System.Windows.Forms.AnchorStyles.Left) 
            | System.Windows.Forms.AnchorStyles.Right)));
            this.groupBoxHead.Controls.Add(this.btnAdmin);
            this.groupBoxHead.Controls.Add(this.btnBanTeam);
            this.groupBoxHead.Controls.Add(this.btnTeamCheck);
            this.groupBoxHead.Controls.Add(this.btnEndCup);
            this.groupBoxHead.Controls.Add(this.btnLogout);
            this.groupBoxHead.Location = new System.Drawing.Point(12, 28);
            this.groupBoxHead.Name = "groupBoxHead";
            this.groupBoxHead.Size = new System.Drawing.Size(1288, 66);
            this.groupBoxHead.TabIndex = 2;
            this.groupBoxHead.TabStop = false;
            // 
            // btnAdmin
            // 
            this.btnAdmin.Font = new System.Drawing.Font("Microsoft Sans Serif", 12F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            this.btnAdmin.Location = new System.Drawing.Point(416, 19);
            this.btnAdmin.Name = "btnAdmin";
            this.btnAdmin.Size = new System.Drawing.Size(89, 30);
            this.btnAdmin.TabIndex = 0;
            this.btnAdmin.Text = "Admin";
            this.btnAdmin.UseVisualStyleBackColor = true;
            this.btnAdmin.Visible = false;
            this.btnAdmin.Click += new System.EventHandler(this.btnAdmin_Click);
            // 
            // btnBanTeam
            // 
            this.btnBanTeam.Font = new System.Drawing.Font("Microsoft Sans Serif", 12F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            this.btnBanTeam.Location = new System.Drawing.Point(139, 19);
            this.btnBanTeam.Name = "btnBanTeam";
            this.btnBanTeam.Size = new System.Drawing.Size(126, 30);
            this.btnBanTeam.TabIndex = 0;
            this.btnBanTeam.Text = "Teams Bannen";
            this.btnBanTeam.UseVisualStyleBackColor = true;
            this.btnBanTeam.Click += new System.EventHandler(this.btnBanTeam_Click);
            // 
            // btnTeamCheck
            // 
            this.btnTeamCheck.Font = new System.Drawing.Font("Microsoft Sans Serif", 12F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            this.btnTeamCheck.Location = new System.Drawing.Point(7, 19);
            this.btnTeamCheck.Name = "btnTeamCheck";
            this.btnTeamCheck.Size = new System.Drawing.Size(126, 30);
            this.btnTeamCheck.TabIndex = 0;
            this.btnTeamCheck.Text = "Teams Prüfen";
            this.btnTeamCheck.UseVisualStyleBackColor = true;
            this.btnTeamCheck.Click += new System.EventHandler(this.btnTeamCheck_Click);
            // 
            // btnEndCup
            // 
            this.btnEndCup.Font = new System.Drawing.Font("Microsoft Sans Serif", 12F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            this.btnEndCup.Location = new System.Drawing.Point(271, 19);
            this.btnEndCup.Name = "btnEndCup";
            this.btnEndCup.Size = new System.Drawing.Size(139, 30);
            this.btnEndCup.TabIndex = 0;
            this.btnEndCup.Text = "Cup Abschließen";
            this.btnEndCup.UseVisualStyleBackColor = true;
            this.btnEndCup.Click += new System.EventHandler(this.btnEndCup_Click);
            // 
            // btnLogout
            // 
            this.btnLogout.Font = new System.Drawing.Font("Microsoft Sans Serif", 12F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            this.btnLogout.Location = new System.Drawing.Point(1193, 19);
            this.btnLogout.Name = "btnLogout";
            this.btnLogout.Size = new System.Drawing.Size(89, 30);
            this.btnLogout.TabIndex = 0;
            this.btnLogout.Text = "Logout";
            this.btnLogout.UseVisualStyleBackColor = true;
            this.btnLogout.Click += new System.EventHandler(this.btnLogout_Click);
            // 
            // groupBoxValue
            // 
            this.groupBoxValue.Anchor = ((System.Windows.Forms.AnchorStyles)((((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Bottom) 
            | System.Windows.Forms.AnchorStyles.Left) 
            | System.Windows.Forms.AnchorStyles.Right)));
            this.groupBoxValue.Controls.Add(this.tabControl1);
            this.groupBoxValue.Location = new System.Drawing.Point(12, 123);
            this.groupBoxValue.Name = "groupBoxValue";
            this.groupBoxValue.Size = new System.Drawing.Size(1288, 647);
            this.groupBoxValue.TabIndex = 2;
            this.groupBoxValue.TabStop = false;
            // 
            // tabControl1
            // 
            this.tabControl1.Controls.Add(this.tabPage1);
            this.tabControl1.Controls.Add(this.tabPage2);
            this.tabControl1.Location = new System.Drawing.Point(7, 20);
            this.tabControl1.Name = "tabControl1";
            this.tabControl1.SelectedIndex = 0;
            this.tabControl1.Size = new System.Drawing.Size(1275, 621);
            this.tabControl1.TabIndex = 0;
            // 
            // tabPage1
            // 
            this.tabPage1.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(45)))), ((int)(((byte)(39)))), ((int)(((byte)(89)))));
            this.tabPage1.Controls.Add(this.dtGrVPolska);
            this.tabPage1.Location = new System.Drawing.Point(4, 22);
            this.tabPage1.Name = "tabPage1";
            this.tabPage1.Padding = new System.Windows.Forms.Padding(3);
            this.tabPage1.Size = new System.Drawing.Size(1267, 595);
            this.tabPage1.TabIndex = 0;
            this.tabPage1.Text = "Polska Cup";
            // 
            // dtGrVPolska
            // 
            this.dtGrVPolska.BackgroundColor = System.Drawing.Color.FromArgb(((int)(((byte)(45)))), ((int)(((byte)(39)))), ((int)(((byte)(89)))));
            this.dtGrVPolska.ColumnHeadersHeightSizeMode = System.Windows.Forms.DataGridViewColumnHeadersHeightSizeMode.AutoSize;
            this.dtGrVPolska.Location = new System.Drawing.Point(3, 6);
            this.dtGrVPolska.Name = "dtGrVPolska";
            this.dtGrVPolska.Size = new System.Drawing.Size(1258, 586);
            this.dtGrVPolska.TabIndex = 0;
            // 
            // tabPage2
            // 
            this.tabPage2.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(45)))), ((int)(((byte)(39)))), ((int)(((byte)(89)))));
            this.tabPage2.Controls.Add(this.dtGrVSkBCup);
            this.tabPage2.Location = new System.Drawing.Point(4, 22);
            this.tabPage2.Name = "tabPage2";
            this.tabPage2.Padding = new System.Windows.Forms.Padding(3);
            this.tabPage2.Size = new System.Drawing.Size(1267, 595);
            this.tabPage2.TabIndex = 1;
            this.tabPage2.Text = "SkinBaron Cup";
            // 
            // dtGrVSkBCup
            // 
            this.dtGrVSkBCup.BackgroundColor = System.Drawing.Color.FromArgb(((int)(((byte)(45)))), ((int)(((byte)(39)))), ((int)(((byte)(89)))));
            this.dtGrVSkBCup.ColumnHeadersHeightSizeMode = System.Windows.Forms.DataGridViewColumnHeadersHeightSizeMode.AutoSize;
            this.dtGrVSkBCup.Location = new System.Drawing.Point(3, 6);
            this.dtGrVSkBCup.Name = "dtGrVSkBCup";
            this.dtGrVSkBCup.Size = new System.Drawing.Size(1259, 586);
            this.dtGrVSkBCup.TabIndex = 1;
            // 
            // Main
            // 
            this.AutoScaleDimensions = new System.Drawing.SizeF(6F, 13F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.BackColor = System.Drawing.Color.FromArgb(((int)(((byte)(45)))), ((int)(((byte)(39)))), ((int)(((byte)(89)))));
            this.ClientSize = new System.Drawing.Size(1312, 782);
            this.Controls.Add(this.groupBoxValue);
            this.Controls.Add(this.groupBoxHead);
            this.Icon = ((System.Drawing.Icon)(resources.GetObject("$this.Icon")));
            this.Name = "Main";
            this.Text = "SkinBaronCup";
            this.groupBoxHead.ResumeLayout(false);
            this.groupBoxValue.ResumeLayout(false);
            this.tabControl1.ResumeLayout(false);
            this.tabPage1.ResumeLayout(false);
            ((System.ComponentModel.ISupportInitialize)(this.dtGrVPolska)).EndInit();
            this.tabPage2.ResumeLayout(false);
            ((System.ComponentModel.ISupportInitialize)(this.dtGrVSkBCup)).EndInit();
            this.ResumeLayout(false);

        }

        #endregion

        private System.Windows.Forms.GroupBox groupBoxHead;
        private System.Windows.Forms.Button btnAdmin;
        private System.Windows.Forms.Button btnLogout;
        private System.Windows.Forms.GroupBox groupBoxValue;
        private System.Windows.Forms.Button btnEndCup;
        private System.Windows.Forms.Button btnBanTeam;
        private System.Windows.Forms.Button btnTeamCheck;
        private System.Windows.Forms.TabControl tabControl1;
        private System.Windows.Forms.TabPage tabPage1;
        private System.Windows.Forms.TabPage tabPage2;
        private System.Windows.Forms.DataGridView dtGrVPolska;
        private System.Windows.Forms.DataGridView dtGrVSkBCup;
    }
}