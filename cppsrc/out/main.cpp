/* Fer un programa que llegeixi una seqüència de strings acabada amb marca de fi "#" i els mostri en ordre invers al que s’ha entrat per teclat. Cal usar la classe PilaString i no hi ha cap limitació sobre la quantitat d’elements a processar. Nota: Suposem que l’entrada són cadenes sense espais.
*/
#include <iostream>
#include <string>
#include "PilaString.h"

using namespace std;

int main() {
    PilaString p;
    
    cout << "ENTRA TEXT ACABAT EN #:" << endl;
    string s;
    cin >> s;
    while (s!="#") {
        p.empila(s);
        cin >> s;
    }
    
    cout << "TEXT INVERTIT:" << endl;
    while (!p.buida()) {
        cout << p.cim() << " ";
        p.desempila();
    }
    cout << endl;
    
    return 0;
}
