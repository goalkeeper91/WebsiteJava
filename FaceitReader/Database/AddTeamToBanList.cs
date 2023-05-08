using MySql.Data.MySqlClient;
using Newtonsoft.Json;
using Newtonsoft.Json.Linq;
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Net;

namespace FaceitReader.Database
{
    internal class AddTeamToBanList
    {
        public static bool addTeam(string teamname, string time, List<Tuple<string,string,string>> list)
        {
            string querySelect = "Select * from banned_teams where teamname = '"+teamname+"'";
            string queryInsert = "INSERT INTO banned_teams (team_id,teamname,team_link,length) VALUES('";
            MySqlConnection cnn = Connection.openConnection();
            MySqlConnection cnn2 = Connection.openConnection();
            MySqlCommand cmd;
            MySqlDataReader rdr;
            bool succesfull = false;
            cmd = new MySqlCommand(querySelect,cnn);
            rdr = cmd.ExecuteReader();
            int length;
            if (time == "permanent")
            {
                length = 5214;
            }
            else
            {
                length = Convert.ToInt32(time);
            }    

            if (!rdr.HasRows)
            {
                if(list.Count() == 0)
                {
                    succesfull = false;
                }
                else
                {
                    succesfull = true;
                    queryInsert = queryInsert + list[0].Item1 + "','" + list[0].Item2 + "','" + list[0].Item3 + "','" + length + "')";
                    cmd = new MySqlCommand(queryInsert,cnn2);
                    cmd.ExecuteNonQuery();
                }
                
            }
            rdr.Close();
            cnn.Close();

            return succesfull;
        }

        public static bool addPlayer(string cupId, string teamId)
        {
            bool succesfull = false;

            for (int i = 0; i <= 140; i += 20)
            {
                string o = i.ToString();
                string url = "https://api.faceit.com/championships/v1/championship/" + cupId + "/subscription?limit=20&offset=" + o + "";
                var httpPRequest = (HttpWebRequest)WebRequest.Create(url);
                var httpPResponse = (HttpWebResponse)httpPRequest.GetResponse();

                using (var streamReader = new StreamReader(httpPResponse.GetResponseStream()))
                {
                    var result = streamReader.ReadToEnd();

                    JObject json = JsonConvert.DeserializeObject<JObject>(result);

                    if (json["payload"]["items"] != null)
                    {
                        foreach (var payload in json["payload"]["items"])
                        {
                            if (payload["team"]["id"].ToString() == teamId)
                            {
                                var id = payload["team"]["id"].ToString();
                                foreach (var member in payload["team"]["members"])
                                {
                                    playerToBanList(member["id"].ToString(), member["nickname"].ToString(), id);
                                }
                                succesfull = true;
                            }
                        }
                    }
                }
            }

            return succesfull;
        }

        public static void playerToBanList(string playerId, string playerName, string teamId)
        {
            string querySelect = "Select * from banned_players where player_id = '" + playerId + "'";
            string queryInsert = "INSERT INTO banned_players (player_id,player_name,team_id) VALUES('";
            MySqlConnection cnn = Connection.openConnection();
            MySqlConnection cnn2 = Connection.openConnection();
            MySqlCommand cmd;
            MySqlDataReader rdr;
            cmd = new MySqlCommand(querySelect, cnn);
            rdr = cmd.ExecuteReader();
            
            if (!rdr.HasRows)
            {
                queryInsert = queryInsert + playerId + "','" + playerName + "','" + teamId + "')";
                cmd = new MySqlCommand(queryInsert, cnn2);
                cmd.ExecuteNonQuery();
            }
            rdr.Close();
            cnn.Close();
        }
    }
}
