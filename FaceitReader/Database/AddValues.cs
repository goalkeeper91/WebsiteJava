using MySql.Data.MySqlClient;
using Newtonsoft.Json;
using Newtonsoft.Json.Linq;
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Net;
using System.Text;
using System.Threading.Tasks;

namespace FaceitReader.Database
{
    internal class AddValues
    {
        public static void fillData(string id, int match)
        {
            var Murl = "https://open.faceit.com/data/v4/championships/" + id + "/matches?type=past&offset=0&limit=" + match;
            var httpMRequest = (HttpWebRequest)WebRequest.Create(Murl);


            httpMRequest.ContentType = "application/json; charset=utf-8";
            httpMRequest.Accept = "application/json";
            httpMRequest.Headers["Authorization"] = Globals.Api;
            var httpMResponse = (HttpWebResponse)httpMRequest.GetResponse();
            MySqlConnection cnn = Connection.openConnection();
            MySqlConnection cnn2 = Connection.openConnection();
            MySqlCommand command;
            MySqlDataReader reader1;
            MySqlDataReader reader2;


            using (var streamReader = new StreamReader(httpMResponse.GetResponseStream()))
            {
                var result = streamReader.ReadToEnd();
                var x = 0;

                JObject json = JsonConvert.DeserializeObject<JObject>(result);

                foreach (JToken item in json["items"])
                {
                    foreach (JToken placement in item["teams"])
                    {
                        foreach (JToken team in placement.Children())
                        {
                            string queryT = "Select * from teams where team_id ='" + team["faction_id"].ToString() + "'";
                            string queryI = "INSERT INTO teams (team_id,teamname) VALUES('"+ team["faction_id"].ToString()+ "','"+ team["name"].ToString()+"')";
                            command = new MySqlCommand(queryT, cnn2);
                            reader2 = command.ExecuteReader();
                            if (!reader2.HasRows)
                            {
                                command = new MySqlCommand(queryI, cnn);
                                command.ExecuteNonQuery();
                            }
                            reader2.Close();

                            if (team["faction_id"].ToString() != "bye")
                            {
                                foreach (JToken player in team["roster"])
                                {
                                    x++;
                                    string query2 = "Select * from player where player_id ='" + player["player_id"].ToString() + "'";
                                    string query = "";
                                    command = new MySqlCommand(query2, cnn2);
                                    reader1 = command.ExecuteReader();
                                    if (!reader1.HasRows)
                                    {
                                        if (team["faction_id"].ToString() != "bye")
                                        {
                                            query = "INSERT INTO player (team_id,player_id,cup_id,playername,game_player_id) VALUES('" + team["faction_id"].ToString() + "','" + player["player_id"].ToString() + "','" + id + "','" + player["nickname"].ToString() + "','" + player["game_player_id"].ToString() + "')";
                                        }
                                    }
                                    else
                                    {
                                        query = "UPDATE player SET cup_id='" + id + "' WHERE player_id = '" + player["player_id"].ToString() + "'";
                                    }
                                    reader1.Close();
                                    command = new MySqlCommand(query, cnn);
                                    command.ExecuteNonQuery();
                                }
                            }
                        }
                    }
                }
                cnn2.Close();
                cnn.Close();
            }
        }
    }
}
