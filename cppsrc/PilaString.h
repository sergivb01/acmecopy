#ifndef PILASTRING_H
#define PILASTRING_H
#include <string>
using namespace std;

class PilaString {
    // Descripcio: una pila d’strings

public:
    // CONSTRUCTORS I DESTRUCTOR ----------------------------------
    PilaString();
    // Pre: --; Post: pila buida
    PilaString(const PilaString& o);  // const. de copia
    // Pre: --; Post: aquesta pila es copia de la pila o
    ~PilaString();
    // Pre: --; Post: memoria alliberada (inclosa la dinàmica)
    
    // CONSULTORS -------------------------------------------------
    bool buida() const;
    // Pre: -- ; Post: retorna cert si la pila es buida; fals en c.c.
    string cim() const;
    // Pre: pila no buida; Post: retorna el valor del cim de la pila
    
    // MODIFICADORS -----------------------------------------------
    void empila(string s);
    // Pre: --; Post: ha afegit s a dalt de la pila
    void desempila();
    // Pre: pila no buida; Post: ha eliminat element de dalt de la pila
    
    // OPERADORS REDEFINITS ---------------------------------------
    PilaString& operator=(const PilaString& o);
    // Pre: -- ; Post: aquesta pila es copia de la pila o i la seva memòria dinàmica ha estat prèviament alliberada

private:
    struct Node {
        string valor;
        Node* seg;
    };
    
    // ATRIBUTS
        Node* a_cim; // punter al cim de la pila

    // METODES PRIVATS
    void copia(const PilaString& o);
    // Pre: pila buida; Post: aquesta pila es copia de la pila o
    void allibera();
    // Pre: --; Post: memoria dinàmica alliberada
};

#endif // PILASTRING_H