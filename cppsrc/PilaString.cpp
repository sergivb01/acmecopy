#include "PilaString.h"
#include <iostream>
#include <string>

PilaString::PilaString(){
    a_cim= NULL;
}

PilaString::PilaString(const PilaString& o){
    a_cim= NULL;
    copia(o); // crida a mètode privat
}

PilaString::~PilaString(){
    allibera(); // crida a mètode privat
}

// CONSULTORS -------------------------------------------------
bool PilaString::buida() const{
    return a_cim==NULL;
}

string PilaString::cim() const{
    return a_cim->valor;
}

// MODIFICADORS -----------------------------------------------
void PilaString::empila(string s){
    Node* nou= new Node; // necessari reservar memoria
    nou->valor= s;
    nou->seg= a_cim;
    a_cim= nou;
}

void PilaString::desempila(){
    Node* aux= a_cim;
    a_cim= a_cim->seg;
    delete aux;
}

// OPERADORS ---------------------------------------------
PilaString& PilaString::operator=(const PilaString& o){
    if (this != &o){
        allibera();
        copia(o);
    }
    return *this;
}

// METODES PRIVATS ------------------------------------------
void PilaString::copia(const PilaString& o) {
    if(!o.buida()){
        Node *p, *q, *aux;
        p = new Node; a_cim = p;
        aux = o.a_cim; p->valor = aux->valor; p->seg = NULL;
        q = p;
        while (aux->seg!=NULL) {
            aux = aux->seg;
            p = new Node; q->seg = p;
            p->valor = aux->valor; p->seg = NULL;
            q = p;
        }
    }
}

void PilaString::allibera(){
    while (!buida()) {
        Node* aux= a_cim;
        a_cim= a_cim->seg;
        delete aux;
    }
}