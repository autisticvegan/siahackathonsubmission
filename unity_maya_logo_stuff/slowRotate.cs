using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class slowRotate : MonoBehaviour
{
    float _rotationSpeed = 600f;
     float _rotationSpeedDecrease = 165.5f;
    // Start is called before the first frame update
    void Start()
    {

       
    }

    // Update is called once per frame
    void Update()
    {
             if (_rotationSpeed >= 0)
             {
                 _rotationSpeed -= _rotationSpeedDecrease * Time.deltaTime;
 
                 transform.Rotate(0, 0, _rotationSpeed * Time.deltaTime);
             }
    }
}
