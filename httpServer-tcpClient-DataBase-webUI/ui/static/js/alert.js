  const getData = async () => {
    try {
        const response = await fetch('alert');
        if (response.status === 200) {
            const data = await response.json(); //extract JSON from the http response
            if (data['IVL'] == true) alert('Пациент откинул тапки!')
            console.log(data);               
        }
    } catch (err) {
        console.log(err);
    } finally {
        setTimeout(getData , 2000);
    }
  };
  getData()