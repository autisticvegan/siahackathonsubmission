using System;
using System.IO;
using System.Security.Cryptography;
using System.Text;
using SiaSkynet;

namespace ConsoleApp1
{
    class Program
    {
        //https://stackoverflow.com/a/1344255
        public static string GetUniqueString(int size)
        {
            char[] chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890".ToCharArray();
           byte[] data = new byte[4 * size];
           using (RNGCryptoServiceProvider crypto = new RNGCryptoServiceProvider())
           {
              crypto.GetBytes(data);
           }
           StringBuilder result = new StringBuilder(size);
           for (int i = 0; i < size; i++)
           {
                var rnd = BitConverter.ToUInt32(data, i * 4);
                var idx = rnd % chars.Length;

                result.Append(chars[idx]);
            }

           return result.ToString();
        }

        static void Main(string[] args)
        {
            if (args.Length != 2)
            {
                Console.WriteLine("Invalid args - should only be the datakey and the file containing the data you want to put");
                return;
            }

            SiaSkynet.SiaSkynetClient c = new SiaSkynetClient();
            bool seedFileExists = File.Exists(@".\seed.txt");
            string seed = "";
            if (seedFileExists)
            {
                Console.WriteLine("seed.txt found, generating keys off it");
                seed = System.IO.File.ReadAllText(@".\seed.txt");
            } else
            {
                string generatedSeed = GetUniqueString(64);
                System.IO.File.WriteAllText(@".\seed.txt", generatedSeed);
                seed = generatedSeed;
                Console.WriteLine("seed.txt not found, generating a seed randomly and writing it to seed.txt");
            }
            var keys = SiaSkynetClient.GenerateKeys(seed).Result;
            Console.WriteLine("Updating registry at " + args[0] + " with content in " + args[1]);
            string content = System.IO.File.ReadAllText(args[1]);
            bool r = c.SkyDbSet(keys.privateKey, keys.publicKey, args[0], content).Result;

            Console.WriteLine("Registry now contains:");
            string actualContent = c.SkyDbGetAsString(keys.publicKey, args[0]).Result;
            Console.WriteLine(actualContent);
        }
    }
}
