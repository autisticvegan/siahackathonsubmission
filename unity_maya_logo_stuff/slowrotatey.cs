using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class slowrotatey : MonoBehaviour
{
    float _rotationSpeed = 600f;
     float _rotationSpeedDecrease = 160f;
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
 
                 transform.Rotate(_rotationSpeed * Time.deltaTime, 0, 0);
             }
    }
}
