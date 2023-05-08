using MySql.Data.MySqlClient;

namespace FaceitReader.Database
{
    internal class AddVoucherToDB
    {
        public static void VoucherAdd(string values ,string cupname, string voucher)
        {
            string query = "Insert Into vouchers (cupname,value,voucher) Values('" + cupname + "','" + values + "','" + voucher + "')";
            MySqlConnection conn = Connection.openConnection();
            MySqlCommand cmd;

            cmd = new MySqlCommand(query, conn);
            cmd.ExecuteNonQuery();
            conn.Close();
        }
    }
}
