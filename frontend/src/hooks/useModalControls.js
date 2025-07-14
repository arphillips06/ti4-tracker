import { useState } from 'react';

export default function useModalControls() {
    const [modals, setModals] = useState({
        crown: false,
        shard: false,
        seed: false,
        imperial: false,
        obsidian: false,
        agenda: false,
        censure: false,
    });

    const toggleModal = (key, value = null) => {
        setModals(prev => ({
            ...prev,
            [key]: value === null ? !prev[key] : value,
        }));
    };

    return {
        modals,
        toggleModal,
    };
}
