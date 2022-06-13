using System;

namespace FaceitReader.Classes
{
    public class FullPlaceList : IEquatable<FullPlaceList>
    {
        public int Id { get; set; }
        public int Place { get; set; }
        public string TeamName { get; set; }
        public string CaptainName { get; set; }
        public string Prize { get; set; }
        public override string ToString()
        {
            return "Place: " + Place + "  TeamName: " + TeamName + "   CaptainName: " + CaptainName + "   Prize: " + Prize;
        }
        public override bool Equals(object obj)
        {
            if (obj == null) return false;
            FullPlaceList objAsFullPlaceList = obj as FullPlaceList;
            if (objAsFullPlaceList == null) return false;
            else return Equals(objAsFullPlaceList);
        }
        public override int GetHashCode()
        {
            return Id;
        }
        public bool Equals(FullPlaceList other)
        {
            if (other == null) return false;
            return (this.Id.Equals(other.Id));
        }
    }
}
