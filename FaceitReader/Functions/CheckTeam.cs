using FaceitReader.Database;
using Newtonsoft.Json;
using Newtonsoft.Json.Linq;
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Net;
using System.Text;
using System.Threading.Tasks;

namespace FaceitReader.Functions
{
    internal class CheckTeam
    {

        public static List<string> isTeamBanned(string cup_id, int limit)
        {
            string url = "https://api.faceit.com/championships/v1/championship/"+cup_id+"/subscription?limit="+limit+"&offset=0";
            string team_id;
            List<string> bannedTeams = new List<string>();

            var httpPRequest = (HttpWebRequest)WebRequest.Create(url);

            httpPRequest.ContentType = "application/json; charset=utf-8";
            httpPRequest.Accept = "application/json";
            var httpPResponse = (HttpWebResponse)httpPRequest.GetResponse();

            using (var streamReader = new StreamReader(httpPResponse.GetResponseStream()))
            {
                var result = streamReader.ReadToEnd();

                JObject json = JsonConvert.DeserializeObject<JObject>(result);

                foreach (var item in json["items"])
                {
                    foreach (var team in item["team"])
                    {
                        team_id = GetBannedTeams.getBannedTeam(team["id"].ToString());
                        if (!string.IsNullOrWhiteSpace(team_id))
                        {
                            if (team["id"].ToString() == team_id)
                            {
                                bannedTeams.Add(team["name"].ToString());
                            }
                        }
                    }
                }
            }

            return bannedTeams;
        }
    }
}
