using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class camera : MonoBehaviour
{    
    float zoomspeed = 2.5f;
     float zoomspeeddecrease = .8f;
    public Transform target;
    // Start is called before the first frame update
    void Start()
    {
        
    }

    // Update is called once per frame
    void Update()
    {
        transform.LookAt(target);
            if (zoomspeed >= 0)
             {
                 zoomspeed -= zoomspeeddecrease * Time.deltaTime;
                Debug.Log("zoomspeed is " + zoomspeed);
                 transform.Translate(0, 0, zoomspeed * Time.deltaTime);
             }
    }
}



